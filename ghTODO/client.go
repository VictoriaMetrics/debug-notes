package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/HdrHistogram/hdrhistogram-go"
	"github.com/VictoriaMetrics/VictoriaMetrics/lib/prompb"
	"github.com/VividCortex/ewma"
	"github.com/golang/snappy"
	"github.com/makasim/backpressure"
	"golang.org/x/time/rate"
)

var idx atomic.Int64

// statusStats tracks HTTP status code counts
type statusStats struct {
	counts map[int]*atomic.Int64
}

func newStatusStats() *statusStats {
	return &statusStats{
		counts: make(map[int]*atomic.Int64),
	}
}

func (s *statusStats) record(statusCode int) {
	if _, exists := s.counts[statusCode]; !exists {
		s.counts[statusCode] = &atomic.Int64{}
	}
	s.counts[statusCode].Add(1)
}

func (s *statusStats) snapshot() map[int]int64 {
	result := make(map[int]int64)
	for code, counter := range s.counts {
		result[code] = counter.Load()
	}
	return result
}

var stats = newStatusStats()

func main() {
	log.Printf("Start time is %s", time.Now().Format(time.RFC3339))
	if len(os.Args) < 2 {
		log.Fatal("invalid number of arguments")
	}

	rateLimit0 := os.Args[1]
	if rateLimit0 == "" {
		log.Fatalf("invalid rateLimit limit: %q", rateLimit0)
	}
	rateLimit, err := strconv.Atoi(rateLimit0)
	if err != nil {
		log.Fatalf("invalid rate limit %q: %v", rateLimit0, err)
	}

	l := rate.NewLimiter(rate.Limit(rateLimit), rateLimit/5)

	c := &http.Client{
		Timeout: time.Minute,
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			MaxIdleConns:          500,
			MaxIdleConnsPerHost:   500,
			IdleConnTimeout:       30 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}

	//go func() {
	//	t := time.NewTicker(time.Minute)
	//	defer t.Stop()
	//
	//	for range t.C {
	//		idx.Add(1)
	//	}
	//}()

	//hdrhistogram.New(0, 100000000)

	// Start status reporter
	go func() {
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()

		var prevSnapshot map[int]int64
		for range ticker.C {
			currentSnapshot := stats.snapshot()
			log.Println("=== HTTP Status Statistics (last 10s) ===")

			if prevSnapshot != nil {
				for code, count := range currentSnapshot {
					prevCount := prevSnapshot[code]
					delta := count - prevCount
					if delta > 0 {
						log.Printf("  Status %d: %d requests", code, delta)
					}
				}
			} else {
				for code, count := range currentSnapshot {
					log.Printf("  Status %d: %d requests", code, count)
				}
			}

			prevSnapshot = currentSnapshot
		}
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)

	limit := rateLimit + int(float64(rateLimit)*0.2)

	log.Printf("starting %d workers with %d rps limit", limit, rateLimit)

	bp, err := backpressure.New(backpressure.Config{
		MaxMax:       1000,
		Max:          10,
		MinMax:       2,
		DecidePeriod: time.Second,
		//ThresholdPercent: 0.02,
		IncreasePercent: 0.02,
		DecreasePercent: 0.1,
	})
	if err != nil {
		log.Panicf("failed to create backpressure: %s", err)
	}

	for i := 0; i < limit; i++ {
		go func() {
			buf := make([]byte, 0, 1<<20) // 1MB

			for {
				bpt, allow := bp.Acquire()
				if !allow {
					//time.Sleep(time.Millisecond * 100)
					//continue
				}

				if err := l.Wait(context.Background()); err != nil {
					log.Panicf("rate limit wait error: %s", err)
				}

				bpt1 := &bpt

				buf = buf[:0]
				handleRequest(c, bpt1, i, buf)

				bp.Release(bpt)
			}
		}()
	}

	<-sigs
}

func handleRequest(c *http.Client, bpt *backpressure.Token, sender int, buf []byte) {
	payload := genPayload(10000, sender, buf)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*55)
	defer cancel()

	//req, _ := http.NewRequest("POST", `http://127.0.0.1:8429/api/v1/write`, bytes.NewReader(payload))
	//req, _ := http.NewRequest("POST", `http://127.0.0.1:8480/insert/0/prometheus/api/v1/write`, bytes.NewReader(payload))
	req, _ := http.NewRequestWithContext(ctx, "POST", `http://127.0.0.1:8427/insert/0/prometheus/api/v1/write`, bytes.NewReader(payload))
	req.SetBasicAuth("foo", "bar")
	req.Header.Set("Content-Type", "application/x-protobuf")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("User-Agent", "aUserAgent")
	req.Header.Set("Content-Encoding", "snappy")
	req.Header.Set("X-Prometheus-Remote-Write-Version", "0.1.0")
	req.Header.Set("Content-Length", fmt.Sprintf("%d", len(payload)))

	start := time.Now()
	totalRequests.Add(1)
	resp, err := c.Do(req)
	if err != nil {
		hMux.Lock()
		_ = h.Current.RecordValue(time.Since(start).Milliseconds())
		hMux.Unlock()

		log.Println("http: do: ", err)
		return
	}

	if resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode == http.StatusServiceUnavailable {
		bpt.Congested = true
	}

	_, _ = io.ReadAll(resp.Body)
	resp.Body.Close()

	hMux.Lock()
	_ = h.Current.RecordValue(time.Since(start).Milliseconds())
	hMux.Unlock()
}

func genPayload(size int, sender int, buf []byte) []byte {
	r := &prompb.WriteRequest{}
	for i := 0; i < size; i++ {
		r.Timeseries = append(r.Timeseries, prompb.TimeSeries{
			Labels: []prompb.Label{
				{Name: "__name__", Value: fmt.Sprintf(`a_metric_case_%d_%d_%d`, i, sender, idx.Load())},
				{Name: "foo", Value: fmt.Sprintf("fooVal_%d", idx.Load())},
				{Name: "bar", Value: fmt.Sprintf("barVal_%d", idx.Load())},
				{Name: "baz", Value: fmt.Sprintf("bazVal_%d", idx.Load())},
			},
			Samples: []prompb.Sample{
				{Value: float64(rand.Int63n(1000)), Timestamp: time.Now().UnixMilli()},
			},
		})
	}

	payload := r.MarshalProtobuf(buf)
	return snappy.Encode(nil, payload)
}

var totalRequests atomic.Int64
var avgRPS = newAvgRPS()

func newAvgRPS() ewma.MovingAverage {
	avgRPS := ewma.NewMovingAverage(60)

	go func() {
		t := time.NewTicker(time.Second)
		defer t.Stop()

		prevTotal := totalRequests.Load()
		for range t.C {
			currTotal := totalRequests.Load()
			avgRPS.Add(float64(currTotal - prevTotal))
			prevTotal = currTotal
		}
	}()

	return avgRPS
}

var hMux sync.Mutex
var h = newHistogram()

func newHistogram() *hdrhistogram.WindowedHistogram {
	h := hdrhistogram.NewWindowed(2, 0, 120000, 3)

	go func() {
		t := time.NewTicker(2 * time.Second)
		for range t.C {
			hMux.Lock()
			h.Rotate()

			mh := h.Merge()

			valToDur := func(v int64) time.Duration {
				return time.Duration(v) * time.Millisecond
			}

			fmt.Fprintf(os.Stdout, "rps=%.1f\tmin=%v\tp50=%v\tp80=%v\tp95=%v\tp99=%v\tmax=%v\n",
				avgRPS.Value(),
				valToDur(mh.Min()), valToDur(mh.ValueAtQuantile(50)), valToDur(mh.ValueAtQuantile(80)),
				valToDur(mh.ValueAtQuantile(95)), valToDur(mh.ValueAtQuantile(99)), valToDur(mh.Max()))
			hMux.Unlock()
		}
	}()

	return h
}
