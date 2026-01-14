package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/VictoriaMetrics/metrics"
)

var (
	demoCounter = metrics.NewCounter(`demo_counter`)
	scrapesMu   sync.Mutex
	scrapes     map[string]int
)

func main() {
	scrapes = make(map[string]int)

	go func() {
		for {
			demoCounter.Inc()
			time.Sleep(time.Second)
		}
	}()

	var fail atomic.Bool
	failT := time.NewTicker(time.Second * 10)

	t := time.NewTimer(time.Millisecond * 1500)

	http.Handle("/metrics", http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		scrapesMu.Lock()

		log.Printf("Scrape #%d from %s", scrapes[r.RemoteAddr], r.RemoteAddr)
		if fail.Load() {
			if _, found := scrapes[r.RemoteAddr]; !found {
				select {
				case <-failT.C:
					fail.Store(false)
				default:
					scrapesMu.Unlock()
					log.Printf("Send 503 for %s", r.RemoteAddr)
					rw.WriteHeader(http.StatusServiceUnavailable)
					return
				}
			}
		}

		scrapes[r.RemoteAddr]++
		cnt := scrapes[r.RemoteAddr]
		scrapesMu.Unlock()

		select {
		case <-t.C:
			t.Reset(time.Millisecond * 1000)
		default:
			if !fail.Load() && cnt%5 == 0 {
				fail.Store(true)
				failT.Stop()
				failT.Reset(time.Second * 10)

				scrapesMu.Lock()
				delete(scrapes, r.RemoteAddr)
				scrapesMu.Unlock()
				log.Printf("Send 503 for %s", r.RemoteAddr)
				rw.WriteHeader(http.StatusServiceUnavailable)
				return
			}
		}

		metrics.WritePrometheus(rw, false)
	}))

	fmt.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}

//func abort(remoteAddr string, rw http.ResponseWriter) {
//	log.Printf("Abort request from %s", remoteAddr)
//	rwh := rw.(http.Hijacker)
//	conn, _, err := rwh.Hijack()
//	if err != nil {
//		log.Printf("ERROR: hijack connection: %v", err)
//	}
//	conn.Close()
//}
