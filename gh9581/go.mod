module anote

go 1.25.0

require (
	github.com/VictoriaMetrics/VictoriaMetrics v1.115.0
	github.com/golang/snappy v1.0.0
)

require github.com/VictoriaMetrics/easyproto v0.1.4 // indirect

replace github.com/VictoriaMetrics/VictoriaMetrics => ../../VictoriaMetrics
