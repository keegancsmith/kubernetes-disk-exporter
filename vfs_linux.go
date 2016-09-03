package main

import "syscall"

// getVFSStats is adapted from cadvisor FS support
func getVFSStats(path string) (*vfsStats, error) {
	var s syscall.Statfs_t
	if err := syscall.Statfs(path, &s); err != nil {
		return nil, err
	}
	return &vfsStats{
		Total:      uint64(s.Frsize) * s.Blocks,
		Free:       uint64(s.Frsize) * s.Bfree,
		Avail:      uint64(s.Frsize) * s.Bavail,
		Inodes:     uint64(s.Files),
		InodesFree: uint64(s.Ffree),
	}, nil
}
