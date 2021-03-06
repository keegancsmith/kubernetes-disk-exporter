package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/docker/docker/pkg/mount"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	root       = flag.String("root", "/", "Specifies where the rootfs is mounted. Without docker it would be /, but usually you need to mount it with a flag like -v /:/rootfs:ro")
	addr       = flag.String("addr", ":9200", "Address on which to expose metrics and web interface.")
	shortLived = flag.Duration("short-lived", 0, "Process will cleanly shutdown after the specified duration. This is a workaround for https://github.com/keegancsmith/kubernetes-disk-exporter/issues/1")

	namespace = "k8snode"
	subsystem = "filesystem"
	labels    = []string{"name"}

	sizeDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "size_bytes"),
		"Filesystem size in bytes.",
		labels, nil,
	)

	freeDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "free_bytes"),
		"Filesystem free space in bytes.",
		labels, nil,
	)

	availDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "avail_bytes"),
		"Filesystem space available to non-root users in bytes.",
		labels, nil,
	)

	fsCollectTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: "collector",
		Name:      "fs_total",
		Help:      "Total number of times stats have been collected from the FS.",
	})
)

func init() {
	prometheus.MustRegister(&collector{})
	prometheus.MustRegister(fsCollectTotal)
}

type vfsStats struct {
	Total      uint64
	Free       uint64
	Avail      uint64
	Inodes     uint64
	InodesFree uint64
}

type collector struct{}

func (c *collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- sizeDesc
	ch <- freeDesc
	ch <- availDesc
}

func (c *collector) Collect(ch chan<- prometheus.Metric) {
	mounts, err := mount.GetMounts()
	if err != nil {
		log.Printf("getmounts failed: %s", err)
		return
	}
	// example mountpath /rootfs/var/lib/kubelet/plugins/kubernetes.io/gce-pd/mounts/gitserver-prod-1
	prefix := filepath.Join(*root, "/var/lib/kubelet/plugins/kubernetes.io/gce-pd/mounts/") + "/"
	for _, m := range mounts {
		path := m.Mountpoint
		if !strings.HasPrefix(path, prefix) {
			continue
		}
		v, err := getVFSStats(path)
		if err != nil {
			log.Printf("%s statfs failed: %s", path, err)
			continue
		}
		name := strings.TrimPrefix(path, prefix)

		ch <- prometheus.MustNewConstMetric(
			sizeDesc, prometheus.GaugeValue,
			float64(v.Total), name,
		)
		ch <- prometheus.MustNewConstMetric(
			freeDesc, prometheus.GaugeValue,
			float64(v.Free), name,
		)
		ch <- prometheus.MustNewConstMetric(
			availDesc, prometheus.GaugeValue,
			float64(v.Avail), name,
		)
	}
	fsCollectTotal.Inc()
}

func main() {
	flag.Parse()

	if *shortLived != 0 {
		time.AfterFunc(*shortLived, func() { os.Exit(0) })
	}

	log.Println("listening on", *addr)
	http.Handle("/metrics", prometheus.Handler())
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal(err)
	}
}
