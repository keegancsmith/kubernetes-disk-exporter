package main

import (
	"flag"
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/docker/docker/pkg/mount"
)

var root = flag.String("root", "/rootfs", "Specifies where the rootfs is mounted. Without docker it would be /, but usually you need to mount it with a flag like -v /:/rootfs:ro")

type vfsStats struct {
	Total      uint64
	Free       uint64
	Avail      uint64
	Inodes     uint64
	InodesFree uint64
}

func main() {
	mounts, err := mount.GetMounts()
	if err != nil {
		log.Fatal(err)
	}
	// example mountpath /rootfs/var/lib/kubelet/plugins/kubernetes.io/gce-pd/mounts/gitserver-prod-1
	prefix := filepath.Join(*root, "/var/lib/kubelet/plugins/kubernetes.io/gce-pd/mounts/")
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
		fmt.Printf("\n%s %#+v\n", path, v)
		fmt.Println(100 * float64(v.Free) / float64(v.Total))
		fmt.Println("total_gb", v.Total/1000000000)
		fmt.Println("free_gb ", v.Free/1000000000)
		fmt.Println("avail_gb", v.Avail/1000000000)
	}
}
