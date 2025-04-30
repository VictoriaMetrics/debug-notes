# Operator

Create cluster:
```
kind create cluster --name quick-start
```

Download operator
```
./01-download-operator.sh
```

Install operator
```
./02-install-operator.sh
```

Run operator checks
```
./03-check-operator.sh
```

# Storage

Create vmsingle manifest
```
./10-create-vmsingle-demo.sh
```

Install vmsingle
```
./11-install-vmsingle-demo.sh
```

Check store is running
```
./12-check-vmsingle-demo.sh
```

Push some metrics
```
# port forward in one terminal
./13-port-forward-vmsingle-demo.sh

# push metrics in another
./14-import-metrics-to-vmsingle-demo.sh
```

Query metrics
```
# port forward in one terminal
./13-port-forward-vmsingle-demo.sh

# query metrics in another, or explore in vmui http://127.0.0.1:8429/vmui
./15-query-metrics-frmo-vmsingle-demo.sh
```

Let's see where the files are stored
```
./16-list-files-vmsingle-demo.sh
```

# Scrape metrics

TODO

# Alerts

TODO


