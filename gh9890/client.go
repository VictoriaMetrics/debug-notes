package main

import (
	"bytes"
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
)

var idx atomic.Int64

func main() {
	if len(os.Args) < 2 {
		log.Fatal("invalid number of arguments")
	}

	workersNum0 := os.Args[1]
	if workersNum0 == "" {
		log.Fatalf("invalid worker number: %q", workersNum0)
	}
	workersNum, err := strconv.Atoi(workersNum0)
	if err != nil {
		log.Fatalf("invalid workers number %q: %v", workersNum0, err)
	}

	c := &http.Client{
		Timeout: time.Minute,
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

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)

	for i := 0; i < workersNum; i++ {
		go func() {
			for {
				handleRequest(c, i)
				time.Sleep(time.Millisecond * 100)
			}
		}()
	}

	<-sigs
}

func handleRequest(c *http.Client, sender int) {
	payload := genPayload(20000, sender)

	//req, _ := http.NewRequest("POST", `http://127.0.0.1:8429/api/v1/write`, bytes.NewReader(payload))
	req, _ := http.NewRequest("POST", `http://127.0.0.1:8480/insert/0/prometheus/api/v1/write`, bytes.NewReader(payload))
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
