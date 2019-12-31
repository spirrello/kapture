# kapture
Run packet captures on Kubernetes pods


## Test ENV vars

```
CAPTURE_PODS_NAMESPACE=kube-system
CAPTURE_PODS=k8s-app=filebeat
KUBECONFIG=config
```
## Test request

```
curl -H "Content-Type: application/json"  -X POST http://localhost:9090/v1/pods -d '{"label":"app=nginx123", "namespace":"dev"}'
```

## Testing kapture-node API

```
curl -H "Content-Type: application/json" -X POST http://localhost:9091/v1/nodeapi -d '[{"name":"nginx-deployment-bb-hsec-78848db5c9-7pmjb","node":"","ip":"10.172.78.226"}]'
```

## Test tcpdump container
```
docker run -it --privileged=true --net=host --rm busybox-tcpdump --name=busybox-tcpdump
```

