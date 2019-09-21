# kcapture
Run packet captures on Kubernetes pods


## Test ENV vars

```
CAPTURE_PODS_NAMESPACE=kube-system
CAPTURE_PODS=k8s-app=filebeat
KUBECONFIG=config
```
## Test request

```
curl -H "Content-Type: application/json"  -X POST http://localhost:9090/v1/pods -d '{"label":"app=nginx123"}'
```