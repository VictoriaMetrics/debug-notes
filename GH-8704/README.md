# Reproducing unexpected EOF on vmagent

Issue: https://github.com/VictoriaMetrics/VictoriaMetrics/pull/8704

1. Checkout VictoriaMetrics repo

```
git clone git@github.com:VictoriaMetrics/VictoriaMetrics.git
git clone git@github.com:makasim/debug-notes.git

cd VictoriaMetrics 
```

2. Build binaries

```
make vmagent
make victoria-metrics
```

3. Run storage:

```
./bin/victoria-metrics 
```

4. Run vmagent

```
./bin/vmagent    -remoteWrite.url=http://0.0.0.0:8428/api/v1/write -maxConcurrentInserts=2 -insert.maxQueueDuration=10s
```

5. Run load script without timeout:

```
go run ../debug-notes/GH-8704/simulate_remote_write.go 0 
```

You should observe a lot of successful writes and some rare warnings like:

```
2025-04-12T13:49:47.551Z	warn	VictoriaMetrics/app/vmagent/main.go:291	remoteAddr: "127.0.0.1:60319"; requestURI: /api/v1/write; cannot read compressed request in 10 seconds: cannot process insert request for 10.000 seconds because 2 concurrent insert requests are executed. Possible solutions: to reduce workload; to increase compute resources at the server; to increase -insert.maxQueueDuration; to increase -maxConcurrentInserts
```

6. Run load script with 1sec timeout:

```
go run ../debug-notes/GH-8704/simulate_remote_write.go 1s 
```

You should observe a lot of successful writes and some rare warnings like:

```
2025-04-12T13:49:47.374Z	warn	VictoriaMetrics/app/vmagent/main.go:291	remoteAddr: "127.0.0.1:60308"; requestURI: /api/v1/write; cannot read compressed request in 10 seconds: unexpected EOF

2025-04-12T13:49:47.551Z	warn	VictoriaMetrics/app/vmagent/main.go:291	remoteAddr: "127.0.0.1:60319"; requestURI: /api/v1/write; cannot read compressed request in 10 seconds: cannot process insert request for 10.000 seconds because 2 concurrent insert requests are executed. Possible solutions: to reduce workload; to increase compute resources at the server; to increase -insert.maxQueueDuration; to increase -maxConcurrentInserts
```
