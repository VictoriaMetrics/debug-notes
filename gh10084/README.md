Start vmauth:

```
../../VictoriaMetrics/bin/vmauth -auth.config=vmauth.yml
```

Start server:
```
go run server.go
```

Start client (concurrent requests):
```
go run client.go
```

Start client2 (sequential requests with pauses via raw TCP):
```
# Default: 10ms pause every 2KB
go run client2.go

# Custom pause duration
go run client2.go -pause=100ms

# Custom target and path
go run client2.go -pause=50ms -addr=127.0.0.1:8427 -path=/foo
```

Scrape metrics from vmauth and check vmauth dashboard in grafana. 