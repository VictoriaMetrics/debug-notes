Start vmauth:

```
../../VictoriaMetrics/bin/vmauth -auth.config=vmauth.yml
```

Start server:
```
go run server.go
```

Start client:
```
go run client.go
```

Scrape metrics from vmauth and check vmauth dashboard in grafana. 