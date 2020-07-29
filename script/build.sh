#!/bin/bash -e
# assuming single implementation in cmd/wp-sqlite
cd cmd/wp-sqlite
# CGO_ENABLED=1 since sqlite requires it
CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o wp
