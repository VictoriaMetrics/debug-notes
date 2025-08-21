package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/VictoriaMetrics/VictoriaMetrics/lib/prompb"
	"github.com/golang/snappy"
)

func main() {
	c := &http.Client{}

	cnt1 := float64(100000)
	cnt2 := float64(100000)

	for {
		handleRequest(c, cnt1)
		cnt1 += 1000
		time.Sleep(3 * time.Second)

		handleRequest(c, cnt2)
		cnt2 += 100
		time.Sleep(28 * time.Second)
	}
}

func handleRequest(c *http.Client, v float64) {
	payload := genPayload(v)

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

func genPayload(v float64) []byte {
	r := &prompb.WriteRequest{}
	r.Timeseries = append(r.Timeseries, prompb.TimeSeries{
		Labels: []prompb.Label{
			{Name: "__name__", Value: `chf.proxy_FlowsIn_value`},
			{Name: "type", Value: "count"},
			{Name: "cid", Value: "cid58"},
			{Name: "proxyid", Value: "proxyid58"},
		},
		Samples: []prompb.Sample{
			{Value: v, Timestamp: time.Now().UnixMilli()},
		},
	})

	payload := r.MarshalProtobuf(nil)
	return snappy.Encode(nil, payload)
}
