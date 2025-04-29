# Operator

Create cluster:
```
kind create cluster --name quick-start
```

Install CRDs:
```
TODO
```

Install operator
```
export VM_VERSION=`basename $(curl -fs -o/dev/null -w %{redirect_url} https://github.com/VictoriaMetrics/operator/releases/latest)`
echo "VM_VERSION=$VM_VERSION"
wget -O victoria-metrics-operator.yaml https://github.com/VictoriaMetrics/operator/releases/download/$VM_VERSION/install-no-webhook.yaml

kubectl apply -f  victoria-metrics-operator.yaml
```

Check operator is running
```
kubectl get pods -n vm -l "control-plane=vm-operator"
```

Check custom resources are created (TODO: move to install CRDs)
```
kubectl api-resources --api-group=operator.victoriametrics.com
```



Print operator version
```
kubectl get pods -n vm -l "control-plane=vm-operator" -o jsonpath='{range .items[*]}{range .spec.containers[?(@.name=="manager")]}{.image}{"\n"}{end}{end}'
```

Print operator configuration options
```
TODO
```

Print VictoriaMetrics versions:
```
TODO
```

# Storage

Install storage
```
kubectl apply -f vmsingle.yaml
```

Check store is running
```
kubectl get pods -n vm -l "app.kubernetes.io/name=vmsingle"
```

Push some metrics
```
curl -i -X POST \
  --url http://127.0.0.1:8429/api/v1/import/prometheus \
  --header 'Content-Type: text/plain' \
  --data 'a_metric{foo="fooVal"} 1
'
```

Query metrics
```
curl -i --url http://127.0.0.1:8429/api/v1/query_range?query=a_metric
```

or explore in vmui http://127.0.0.1:8429/vmui



