# kubernetes-disk-exporter

Export kubernetes persistent volume metrics to prometheus.

Alternatives:
* `node_exporter` is heavy-weight and doesn't parse k8s metadata.
* `cadvisor` doesn't seem to include the stats I want.

The intention of this is to be a simple daemonset you can deploy, and then
prometheus will tell you size, free, available of each persistent volume you
mount in k8s.

This has yet to be deployed, but you can experiment with it by running

```
$ sudo docker run -v /:/rootfs:ro --rm=true -p 9200:9200 keegancsmith/kubernetes-disk-exporter
$ curl http://localhost:9200/metrics
```

on a kubernetes node. If you have any pods which have mounted `gce-pd` disks,
you should see metrics for them.
