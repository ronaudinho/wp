#!/bin/bash -e
# test everything but internal
# somehow internal still gets into cover.out
# TODO investigate why
go test -v ./... -coverprofile=./cover.out
# test internal dir
cd internal; go test -v ./...; cd ..
# visualize (uncomment the line below)
# go tool cover -html=cover.out
