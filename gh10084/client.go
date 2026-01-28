package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Go(func() {
			for {
				b := make([]byte, 17000)
				for i := 0; i < len(b); i++ {
					b[i] = byte(rand.Int31())
				}
				
				req, err := http.NewRequest("POST", `http://127.0.0.1:8427/foo`, bytes.NewBuffer(b))
				if err != nil {
					log.Fatal(err)
				}

				req.Header.Set("Content-Type", "application/octet-stream")
				req.Header.Set("Authorization", "Basic Zm9vOmJhcg==")

				resp, err := http.DefaultClient.Do(req)
				if err != nil {
					log.Printf("Request failed: %v", err)
					time.Sleep(time.Second)
					continue
				}
				if _, err := io.ReadAll(resp.Body); err != nil {
					log.Printf("Failed to read response: %v", err)
					time.Sleep(time.Second)
					continue
				}
				resp.Body.Close()

				fmt.Println("\nRequest completed successfully!")
				time.Sleep(time.Second)
			}
		})
	}

	wg.Wait()
}
