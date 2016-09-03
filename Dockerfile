FROM scratch
ADD kubernetes-disk-exporter /
ENTRYPOINT ["/kubernetes-disk-exporter"]
