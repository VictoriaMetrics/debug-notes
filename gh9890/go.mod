module backpressure

go 1.25.5

require (
	github.com/VictoriaMetrics/VictoriaMetrics v1.126.0
	github.com/golang/snappy v1.0.0
	golang.org/x/time v0.14.0
)

require github.com/VictoriaMetrics/easyproto v1.0.0 // indirect

replace github.com/VictoriaMetrics/VictoriaMetrics => ../../VictoriaMetrics
