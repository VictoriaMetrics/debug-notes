package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/VictoriaMetrics/VictoriaMetrics/lib/prompbmarshal"
	"github.com/golang/snappy"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("invalid number of arguments")
	}
	timeoutArg := os.Args[1]
	if timeoutArg == "" {
		log.Fatalf("invalid timeout %q", timeoutArg)
	}
	timeout, err := time.ParseDuration(timeoutArg)
	if err != nil {
		log.Fatalf("invalid timeout %q: %v", timeoutArg, err)
	}

	var wg sync.WaitGroup

	concurrencyCh := make(chan struct{}, 300)

	c := &http.Client{
		Timeout: timeout,
	}

	for i := 0; i < 300; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for {
				concurrencyCh <- struct{}{}
				handleRequest(c)
				<-concurrencyCh
			}
		}()
	}

	wg.Wait()
}

func handleRequest(c *http.Client) {
	payload := genPayload(100000)

	req, _ := http.NewRequest("POST", `http://127.0.0.1:8429/api/v1/write`, bytes.NewReader(payload))
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

func genPayload(size int) []byte {
	r := &prompbmarshal.WriteRequest{}
	for i := 0; i < size; i++ {
		r.Timeseries = append(r.Timeseries, prompbmarshal.TimeSeries{
			Labels: []prompbmarshal.Label{
				{Name: "__name__", Value: fmt.Sprintf(`a_metric_case_%d`, i)},
				{Name: "foo", Value: "fooVal"},
				{Name: "bar", Value: "barVal"},
				{Name: "baz", Value: "bazVal"},
			},
			Samples: []prompbmarshal.Sample{
				{Value: 1, Timestamp: time.Now().UnixMilli()},
			},
		})
	}

	payload := r.MarshalProtobuf(nil)
	return snappy.Encode(nil, payload)
}
