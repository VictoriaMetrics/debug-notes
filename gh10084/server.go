package main

import (
	"io"
	"log"
	"math/rand"
	"net/http"
)

var response100KB = make([]byte, 100*1024)

func handler(w http.ResponseWriter, r *http.Request) {
	// Discard request body
	io.Copy(io.Discard, r.Body)
	r.Body.Close()

	// Send 100KB response
	w.Header().Set("Content-Type", "application/octet-stream")
	w.WriteHeader(http.StatusOK)
	w.Write(response100KB)
}

func main() {
	for i := 0; i < len(response100KB); i++ {
		response100KB[i] = byte(rand.Int31())
	}

	http.HandleFunc("/", handler)
	log.Println("Server starting on :8899")
	if err := http.ListenAndServe(":8899", nil); err != nil {
		log.Fatal(err)
	}
}
