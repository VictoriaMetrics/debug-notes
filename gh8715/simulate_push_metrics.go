package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run 8715/main.go <case_id>")
	}
	caseID := os.Args[1]
	if caseID == "" {
		log.Fatal("Usage: go run 8715/main.go <case_id>")
	}

	t := time.NewTicker(time.Second * 2)
	defer t.Stop()

	var cnt int
	for {
		select {
		case <-t.C:
			cnt++
			body := bytes.NewBufferString(fmt.Sprintf(`a_metric_case_%s{foo="fooVal",bar="barVal",baz="bazVal"} %d`, caseID, cnt))

			req, err := http.NewRequest("POST", "http://127.0.0.1:8429/api/v1/import/prometheus", body)
			if err != nil {
				log.Fatal("Error creating request:", err)
			}
			req.Header.Set("Content-Type", "text/plain")

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				log.Fatal("Error making request:", err)
			}

			if resp.StatusCode != http.StatusNoContent {
				log.Fatalf("Expected status code 204, got %d", resp.StatusCode)
			}
		}
	}
}
