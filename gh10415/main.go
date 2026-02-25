package main

import (
	"math"
	"math/rand"
	"net/http"
	"time"

	"github.com/VictoriaMetrics/metrics"
)

var (
	// hLowSkewed: most observations cluster at low values (e.g. fast requests dominate)
	hLowSkewed = metrics.NewHistogram(`test_histogram{distribution="low_skewed"}`)

	// hNormal: bell-curve distribution centered around the midpoint
	hNormal = metrics.NewHistogram(`test_histogram{distribution="normal"}`)

	// hHighSkewed: most observations cluster at high values (mirror of low_skewed)
	hHighSkewed = metrics.NewHistogram(`test_histogram{distribution="high_skewed"}`)
)

func main() {
	go generateMetrics()

	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		metrics.WritePrometheus(w, true)
	})

	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}

func generateMetrics() {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	for {
		// Low-skewed: raise uniform to a high power so most values land near 0.
		// Distribution function: f(x) ~ x^(n-1), heavy weight at low end.
		low := math.Pow(rng.Float64(), 4)
		hLowSkewed.Update(low)

		// Normal: Box-Muller transform, clamped to (0,1), centered at 0.5.
		u1 := rng.Float64()
		u2 := rng.Float64()
		z := math.Sqrt(-2*math.Log(u1)) * math.Cos(2*math.Pi*u2)
		normal := 0.5 + z*0.12
		if normal > 0 && normal < 1 {
			hNormal.Update(normal)
		}

		// High-skewed: mirror of low-skewed — raise to power then flip.
		// Most values land near 1.
		high := 1 - math.Pow(rng.Float64(), 4)
		hHighSkewed.Update(high)

		time.Sleep(time.Millisecond * 100)
	}
}
