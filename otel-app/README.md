# OTEL app

This folder contains a Go app instrumented with OTEL SDK. The app pushes its metrics to OTEL collector (see `collector`).
The OTEL collector can forward metrics via OTLP or Prometheus RemoteWrite elsewhere. To VM, for example.

To start app:
```
go run .
```

Once started, app should be available at [http://localhost:8080/articles](http://localhost:8080/articles). Visiting
this page will update some internal metrics.

App is configured to push its metrics to OTEL collector at `localhost:4317`. 
Download [collector binary](https://opentelemetry.io/docs/collector/installation/) and start it:
```
./collector/otelcol --config=collector/conf.yaml
```

`collector/conf.yaml` contains configuration for collecting and pushing metrics elsewhere. Tweak it to send data
via Otel or Prometheus Remote Write.