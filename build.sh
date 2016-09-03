#!/bin/bash

set -ex

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build .
docker build -t keegancsmith/kubernetes-disk-exporter .
