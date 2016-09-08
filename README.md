# kubernetes-disk-exporter

Export kubernetes persistent volume metrics to prometheus.

Alternatives:
* `node_exporter` is heavy-weight and doesn't parse k8s metadata.
* `cadvisor` doesn't seem to include the stats I want.

## Setup on Kubernetes

If you have prometheus setup with kubernetes service discovery you should be
able to just run:

```
$ kubectl apply -f disk-exporter.yaml
```

This will setup a daemonset running the exporter, as well as a service such
that the service discovery can find the prometheus metric endpoints. If you
are not using the example prometheus + k8s config, adjust disk-exporter.yaml
until it works for you.

## Metrics

For each mount the following is exported

* `k8snode_filesystem_size`
* `k8snode_filesystem_avail`
* `k8snode_filesystem_free`

with the labels

* `name`

### Example Queries

The 3 mounts with the least amount of available space as a percentage of total
space.

```
drop_common_labels(bottomk(3, k8snode_filesystem_avail / k8snode_filesystem_size))
```

## Future

* Support more than just `gce-pd` mounts.
* Extract more k8s labels.
* Export more statistics (inodes).

## Testing

Run this on a kubernetes node to see some example output:

```
$ sudo docker run -v /:/rootfs:ro --rm=true -p 9200:9200 keegancsmith/kubernetes-disk-exporter
$ curl http://localhost:9200/metrics
```
