module backpressure

go 1.25.3

require (
	github.com/HdrHistogram/hdrhistogram-go v1.1.2
	github.com/VictoriaMetrics/VictoriaMetrics v1.126.0
	github.com/VividCortex/ewma v1.2.0
	github.com/golang/snappy v1.0.0
	github.com/makasim/backpressure v0.0.0-20251019160754-a2c53b0ef8b0
	golang.org/x/time v0.14.0
)

require github.com/VictoriaMetrics/easyproto v0.1.4 // indirect

replace github.com/VictoriaMetrics/VictoriaMetrics => ../../VictoriaMetrics
