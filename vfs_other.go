// +build !linux

package main

import "errors"

func getVFSStats(path string) (*vfsStats, error) {
	return nil, errors.New("getVFSStats is not implemented for non-linux systems")
}
