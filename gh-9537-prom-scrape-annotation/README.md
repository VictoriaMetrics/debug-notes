
Prepare:
```
kind create cluster --name=promannotation

helm repo add vm https://victoriametrics.github.io/helm-charts/
helm repo update

kubectl create ns vmks
helm show values vm/victoria-metrics-k8s-stack > values.yaml
```

Install VictoriaMetrics stack:
```
helm install vmks vm/victoria-metrics-k8s-stack -f values.yaml -n vmks --debug
```

Install demo app:
```
kubectl apply -f demo-app.yaml
```

Add auto discovery for prometheus.io config:
```

```

