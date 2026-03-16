package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
)

func main() {
	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		dump, err := httputil.DumpRequest(r, true)
		if err != nil {
			panic(fmt.Sprintf("error dumping request: %v", err))
		}
		fmt.Printf("--- request ---\n%s\n", dump)

		fmt.Fprintf(w, "requests_total 123\n")
	})

	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
