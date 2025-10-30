package main

import (
	"math/rand/v2"
	"net/url"
	"sync/atomic"
)

type backendURL struct {
	brokenDeadline     atomic.Uint64
	concurrentRequests atomic.Int32

	url *url.URL
}

func (bu *backendURL) isBroken() bool {
	return false
}

func (bu *backendURL) setBroken() {
	panic("should not be called")
}

func (bu *backendURL) get() {
	bu.concurrentRequests.Add(1)
}

func (bu *backendURL) put() {
	bu.concurrentRequests.Add(-1)
}

// getFirstAvailableBackendURL returns the first available backendURL, which isn't broken.
//
// backendURL.put() must be called on the returned backendURL after the request is complete.
func getFirstAvailableBackendURL(bus []*backendURL) *backendURL {
	bu := bus[0]
	if !bu.isBroken() {
		// Fast path - send the request to the first url.
		bu.get()
		return bu
	}

	// Slow path - the first url is temporarily unavailable. Fall back to the remaining urls.
	for i := 1; i < len(bus); i++ {
		if !bus[i].isBroken() {
			bu = bus[i]
			break
		}
	}
	bu.get()
	return bu
}

// getRoundRobinBackendURL returns the next backendURL in round-robin order, which isn't broken.
//
// backendURL.put() must be called on the returned backendURL after the request is complete.
func getRoundRobinBackendURL(bus []*backendURL, atomicCounter *atomic.Uint32) *backendURL {
	if len(bus) == 1 {
		// Fast path - return the only backend url.
		bu := bus[0]
		bu.get()
		return bu
	}

	// Round-robin through backends
	n := atomicCounter.Add(1) - 1
	idx := n % uint32(len(bus))
	bu := bus[idx]

	if !bu.isBroken() {
		// Fast path - the selected backend is available
		bu.get()
		return bu
	}

	// Slow path - find the next available backend
	for i := uint32(1); i < uint32(len(bus)); i++ {
		idx = (n + i) % uint32(len(bus))
		bu = bus[idx]
		if !bu.isBroken() {
			bu.get()
			return bu
		}
	}

	// All backends are broken - return the originally selected one
	bu = bus[n%uint32(len(bus))]
	bu.get()
	return bu
}

// getPowerOfTwoRandomBackendURL implements the "power of two choices" load balancing algorithm.
// It randomly selects two backends and returns the one with fewer concurrent requests.
//
// backendURL.put() must be called on the returned backendURL after the request is complete.
func getPowerOfTwoRandomBackendURL(bus []*backendURL, atomicCounter *atomic.Uint32) *backendURL {
	if len(bus) == 1 {
		// Fast path - return the only backend url.
		bu := bus[0]
		bu.get()
		return bu
	}

	// Select two random backends
	idx1 := rand.IntN(len(bus))
	idx2 := rand.IntN(len(bus))

	// Ensure we pick two different backends if possible
	if idx1 == idx2 && len(bus) > 1 {
		idx2 = (idx2 + 1) % len(bus)
	}

	bu1 := bus[idx1]
	bu2 := bus[idx2]

	// If one is broken, prefer the non-broken one
	if bu1.isBroken() && !bu2.isBroken() {
		bu2.get()
		return bu2
	}
	if bu2.isBroken() && !bu1.isBroken() {
		bu1.get()
		return bu1
	}

	// Both are available or both are broken - choose the one with fewer requests
	if bu1.concurrentRequests.Load() <= bu2.concurrentRequests.Load() {
		bu1.get()
		return bu1
	}
	bu2.get()
	return bu2
}

// getLeastLoadedBackendURL returns the backendURL with the minimum number of concurrent requests.
//
// backendURL.put() must be called on the returned backendURL after the request is complete.
func getLeastLoadedBackendURL(bus []*backendURL, atomicCounter *atomic.Uint32) *backendURL {
	if len(bus) == 1 {
		// Fast path - return the only backend url.
		bu := bus[0]
		bu.get()
		return bu
	}

	// Slow path - select other backend urls.
	n := atomicCounter.Add(1) - 1
	for i := uint32(0); i < uint32(len(bus)); i++ {
		idx := (n + i) % uint32(len(bus))
		bu := bus[idx]
		if bu.isBroken() {
			continue
		}
		if bu.concurrentRequests.Load() == 0 {
			// Fast path - return the backend with zero concurrently executed requests.
			// Do not use CompareAndSwap() instead of Load(), since it is much slower on systems with many CPU cores.
			bu.concurrentRequests.Add(1)
			return bu
		}
	}

	// Slow path - return the backend with the minimum number of concurrently executed requests.
	buMin := bus[n%uint32(len(bus))]
	minRequests := buMin.concurrentRequests.Load()
	for _, bu := range bus {
		if bu.isBroken() {
			continue
		}
		if n := bu.concurrentRequests.Load(); n < minRequests || buMin.isBroken() {
			buMin = bu
			minRequests = n
		}
	}
	buMin.get()
	return buMin
}

// getEnhancedLeastLoadedBackendURL returns the backendURL with the minimum number of concurrent requests.
//
// backendURL.put() must be called on the returned backendURL after the request is complete.
func getEnhancedLeastLoadedBackendURL(bus []*backendURL, atomicCounter *atomic.Uint32) *backendURL {
	if len(bus) == 1 {
		// Fast path - return the only backend url.
		bu := bus[0]
		bu.get()
		return bu
	}

	// Slow path - select other backend urls.
	n := atomicCounter.Add(1) - 1
	for i := uint32(0); i < uint32(len(bus)); i++ {
		idx := (n + i) % uint32(len(bus))
		bu := bus[idx]
		if bu.isBroken() {
			continue
		}

		// The Load() in front of CompareAndSwap() avoids CAS overhead for items with values bigger than 0.
		if bu.concurrentRequests.Load() == 0 && bu.concurrentRequests.CompareAndSwap(0, 1) {
			atomicCounter.CompareAndSwap(n+1, idx+1)
			// There is no need in the call bu.get(), because we already incremented bu.concrrentRequests above.
			return bu
		}
	}

	// Slow path - return the backend with the minimum number of concurrently executed requests.
	buMinIdx := n % uint32(len(bus))
	minRequests := bus[buMinIdx].concurrentRequests.Load()
	for i := uint32(0); i < uint32(len(bus)); i++ {
		idx := (n + i) % uint32(len(bus))
		bu := bus[idx]
		if bu.isBroken() {
			continue
		}

		reqs := bu.concurrentRequests.Load()
		if reqs < minRequests || bus[buMinIdx].isBroken() {
			buMinIdx = idx
			minRequests = reqs
		}
	}
	buMin := bus[buMinIdx]
	buMin.get()
	atomicCounter.CompareAndSwap(n+1, buMinIdx+1)
	return buMin
}
