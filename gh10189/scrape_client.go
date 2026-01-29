package main

import (
	"fmt"
	"net/http"

	"github.com/VictoriaMetrics/metrics"
)

var m1 = metrics.NewCounter(`gh10189{anotherl="m1",alabel="foo"}`)
var m2 = metrics.NewCounter(`gh10189{anotherl="m2",alabel="foo"}`)
var m3 = metrics.NewCounter(`gh10189{anotherl="m3",alabel="foo"}`)
var m4 = metrics.NewCounter(`gh10189{anotherl="m1",alabel="bar"}`)
var m5 = metrics.NewCounter(`gh10189{anotherl="m2",alabel="bar"}`)
var m6 = metrics.NewCounter(`gh10189{anotherl="m3",alabel="bar"}`)
var m7 = metrics.NewCounter(`gh10189{anotherl="m1",alabel="baz"}`)
var m8 = metrics.NewCounter(`gh10189{anotherl="m2",alabel="baz"}`)
var m9 = metrics.NewCounter(`gh10189{anotherl="m3",alabel="baz"}`)

var m10 = metrics.NewCounter(`gh10189{anotherl="m4",alabel="foo"}`)

func main() {
	m1.Set(1)
	m2.Set(1)
	m3.Set(1)
	m4.Set(1)
	m5.Set(1)
	m6.Set(1)
	m7.Set(1)
	m8.Set(1)
	m9.Set(1)
	m10.Set(2)

	http.Handle("/metrics", http.HandlerFunc(func(rw http.ResponseWriter, _ *http.Request) {
		metrics.WritePrometheus(rw, false)
	}))

	fmt.Println("Starting server on :8080")
	http.ListenAndServe(":8080", nil)
}
