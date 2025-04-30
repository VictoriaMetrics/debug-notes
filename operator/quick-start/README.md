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
./03-check-operator-is-running.sh
./04-print-vm-resources.sh
./05-print-operator-version.sh
./06-print-operator-default-envs.sh
```

# Storage

Create vmsingle manifest
```
./10-create-vmsingle-manifest.sh
```

Install vmsingle
```
./11-install-vmsingle.sh
```

Check store is running
```
./12-check-vmsingle-is-running.sh
```

Push some metrics
```
# port forward in one terminal
./13-port-forward-vmsingle.sh

# push metrics in another
./14-import-metrics.sh
```

Query metrics
```
# port forward in one terminal
./13-port-forward-vmsingle.sh

# query metrics in another, or explore in vmui http://127.0.0.1:8429/vmui
./15-query-metrics.sh
```

Let's see where the files are stored
```
./16-list-storage-files.sh
```

# Scrape metrics

TODO

# Alerts

TODO


