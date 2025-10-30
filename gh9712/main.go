package main

import (
	"flag"
	"fmt"
	"math"
	"math/rand"
	"net/url"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	// Define CLI flags
	strategy := flag.String("strategy", "least-loaded", "Backend selection strategy: 'first-available', 'round-robin', 'power-of-two', 'least-loaded' or 'enha-least-loaded'")
	flag.Parse()

	const (
		numBackends = 10
		concurrency = 1000
		maxRequests = 5_000_000
		maxLatency  = 100
	)

	// Initialize backend URLs
	backends := make([]*backendURL, numBackends)
	for i := 0; i < numBackends; i++ {
		backends[i] = &backendURL{
			url: &url.URL{
				Scheme: "http",
				Host:   fmt.Sprintf("backend-%d:8080", i),
			},
		}
	}

	// Counter for getLeastLoadedBackendURL
	var atomicCounter atomic.Uint32

	// Track durationDistribution of requests to each backend
	durationDistribution := make([]atomic.Int64, numBackends)

	// Create a map to quickly find backend index
	backendIndex := make(map[*backendURL]int)
	for i, bu := range backends {
		backendIndex[bu] = i
	}

	// Semaphore to control concurrency
	sem := make(chan struct{}, concurrency)
	var wg sync.WaitGroup

	fmt.Printf("Starting test with %d backends, %d concurrent workers, %d total requests\n", numBackends, concurrency, maxRequests)
	fmt.Printf("Request latency picked at random from 1 to %d microseconds\n", maxLatency)
	fmt.Printf("Strategy: %s\n\n", *strategy)

	startTime := time.Now()

	// Generate requests
	for i := 0; i < maxRequests; i++ {
		wg.Add(1)
		sem <- struct{}{} // Acquire semaphore

		go func() {
			defer wg.Done()
			defer func() { <-sem }() // Release semaphore

			// Get backend based on selected strategy
			var bu *backendURL
			switch *strategy {
			case "first-available":
				bu = getFirstAvailableBackendURL(backends)
			case "round-robin":
				bu = getRoundRobinBackendURL(backends, &atomicCounter)
			case "power-of-two":
				bu = getPowerOfTwoRandomBackendURL(backends, &atomicCounter)
			case "least-loaded":
				bu = getLeastLoadedBackendURL(backends, &atomicCounter)
			case "enha-least-loaded":
				bu = getEnhancedLeastLoadedBackendURL(backends, &atomicCounter)
			default:
				panic(fmt.Sprintf("unknown strategy: %s", *strategy))
			}

			// Track which backend was selected
			idx := backendIndex[bu]

			// Add random delay from 1ms to 10ms
			requestLatency := time.Duration(1+rand.Intn(maxLatency)) * time.Microsecond

			// Simulate request processing
			time.Sleep(requestLatency)

			durationDistribution[idx].Add(int64(requestLatency))

			// Release the backend
			bu.put()
		}()
	}

	wg.Wait()
	duration := time.Since(startTime)

	fmt.Printf("Test completed in %v\n\n", duration)

	fmt.Println("Distribution of requests duration across backends:")
	fmt.Println("==========================================")

	// Collect results
	countsDur := make([]int64, numBackends)
	var totalDur int64
	for i := 0; i < numBackends; i++ {
		countsDur[i] = durationDistribution[i].Load()
		totalDur += countsDur[i]
	}
	for i := 0; i < numBackends; i++ {
		fmt.Printf("Backend %d: %s requests (%.2f%%)\n", i, time.Duration(countsDur[i]), float64(countsDur[i])*100/float64(totalDur))
	}

	//Calculate mean
	mean := float64(totalDur) / float64(numBackends)

	// Calculate standard deviation
	var sumSquaredDiff float64
	for i := 0; i < numBackends; i++ {
		diff := float64(countsDur[i]) - mean
		sumSquaredDiff += diff * diff
	}
	stdDev := math.Sqrt(sumSquaredDiff / float64(numBackends))

	fmt.Println("\nStatistics:")
	fmt.Println("===========")
	fmt.Printf("Total requests: %d\n", totalDur)
	fmt.Printf("Mean per backend: %.2f\n", mean)
	fmt.Printf("Standard deviation: %.2f\n", stdDev)
	fmt.Printf("Coefficient of variation: %.2f%%\n", (stdDev/mean)*100)

	// Find min and max
	minCount, maxCount := countsDur[0], countsDur[0]
	for i := 1; i < numBackends; i++ {
		if countsDur[i] < minCount {
			minCount = countsDur[i]
		}
		if countsDur[i] > maxCount {
			maxCount = countsDur[i]
		}
	}
	fmt.Printf("Min requests: %d\n", minCount)
	fmt.Printf("Max requests: %d\n", maxCount)
	fmt.Printf("Range: %d\n", maxCount-minCount)
}
