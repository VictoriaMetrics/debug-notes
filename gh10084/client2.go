package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"time"
)

var (
	pauseDuration = flag.Duration("pause", 10*time.Millisecond, "Duration to pause every 4KB when writing request")
	targetAddr    = flag.String("addr", "127.0.0.1:8427", "Target address")
	requestPath   = flag.String("path", "/foo", "Request path")
)

func main() {
	flag.Parse()

	log.Printf("Starting client with write pause: %v", *pauseDuration)
	log.Printf("Target: %s%s", *targetAddr, *requestPath)

	// Prepare 100KB payload
	payload := make([]byte, 100*1024)
	for i := 0; i < len(payload); i++ {
		payload[i] = byte(rand.Int31())
	}

	requestNum := 0
	for {
		requestNum++
		log.Printf("Sending request #%d", requestNum)

		if err := sendRequest(payload); err != nil {
			log.Printf("Request #%d failed: %v", requestNum, err)
			time.Sleep(time.Second)
			continue
		}

		log.Printf("Request #%d completed successfully!", requestNum)
		time.Sleep(time.Second)
	}
}

func sendRequest(payload []byte) error {
	// Connect to target
	conn, err := net.Dial("tcp", *targetAddr)
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	defer conn.Close()

	// Set deadline for the entire operation
	conn.SetDeadline(time.Now().Add(60 * time.Second))

	// Craft HTTP request
	httpRequest := fmt.Sprintf(
		"POST %s HTTP/1.1\r\n"+
			"Host: %s\r\n"+
			"Content-Type: application/octet-stream\r\n"+
			"Content-Length: %d\r\n"+
			"Connection: close\r\n"+
			"\r\n",
		*requestPath,
		*targetAddr,
		len(payload),
	)

	// Combine headers and payload
	fullRequest := append([]byte(httpRequest), payload...)

	// Write first byte
	if _, err := conn.Write(fullRequest[:4097]); err != nil {
		return fmt.Errorf("failed to write first byte: %w", err)
	}
	log.Printf("Sent first byte, pausing for %v...", *pauseDuration)

	// Pause
	time.Sleep(*pauseDuration)

	// Write the rest
	if _, err := conn.Write(fullRequest[4097:]); err != nil {
		return fmt.Errorf("failed to write remaining bytes: %w", err)
	}
	log.Printf("Finished writing %d bytes total, now reading response...", len(fullRequest))

	b, err := io.ReadAll(conn)
	if err != nil {
		return fmt.Errorf("failed to read full response: %w", err)
	}

	log.Printf("Received %d bytes in response", len(b)+1)
	return nil
}
