module backpressure

go 1.25.1

require (
	github.com/HdrHistogram/hdrhistogram-go v1.1.2
	github.com/VictoriaMetrics/VictoriaMetrics v1.126.0
	github.com/golang/snappy v1.0.0
	golang.org/x/time v0.12.0
)

require github.com/VictoriaMetrics/easyproto v0.1.4 // indirect

replace github.com/VictoriaMetrics/VictoriaMetrics => ../../VictoriaMetrics
