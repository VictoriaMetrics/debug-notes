package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync/atomic"
	"time"

	//"github.com/HdrHistogram/hdrhistogram-go"
	"github.com/VictoriaMetrics/VictoriaMetrics/lib/prompb"
	"github.com/golang/snappy"
	"golang.org/x/time/rate"
)

var idx atomic.Int64

func main() {
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

	l := rate.NewLimiter(rate.Limit(rateLimit), rateLimit)

	c := &http.Client{
		Timeout: time.Minute,
	}

	go func() {
		t := time.NewTicker(time.Minute)
		defer t.Stop()

		for range t.C {
			idx.Add(1)
		}
	}()

	//hdrhistogram.New(0, 100000000)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)

	for i := 0; i < 900; i++ {
		go func() {
			for {
				if err := l.Wait(context.Background()); err != nil {
					log.Panicf("rate limit wait error: %s", err)
				}

				handleRequest(c, i%10)
			}
		}()
	}

	<-sigs
}

func handleRequest(c *http.Client, sender int) {
	payload := genPayload(10000, sender)

	req, _ := http.NewRequest("POST", `http://127.0.0.1:8429/api/v1/write`, bytes.NewReader(payload))
	//req, _ := http.NewRequest("POST", `http://127.0.0.1:8480/insert/0/prometheus/api/v1/write`, bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/x-protobuf")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("User-Agent", "aUserAgent")
	req.Header.Set("Content-Encoding", "snappy")
	req.Header.Set("X-Prometheus-Remote-Write-Version", "0.1.0")
	req.Header.Set("Content-Length", fmt.Sprintf("%d", len(payload)))

	resp, err := c.Do(req)
	if err != nil {
		log.Println("http: do: ", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 204 {
		log.Println("http: do: ", resp.Status)
	}
}

func genPayload(size int, sender int) []byte {
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
				{Value: 1, Timestamp: time.Now().UnixMilli()},
			},
		})
	}

	payload := r.MarshalProtobuf(nil)
	return snappy.Encode(nil, payload)
}
