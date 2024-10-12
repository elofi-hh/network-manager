#!/usr/bin/env bash -e
CGO_ENABLED=0 GOOS=linux GOARCH=mips GOMIPS=softfloat go build -trimpath -ldflags="-s -w -extldflags=-static" .
scp manager root@192.168.2.1:/root
