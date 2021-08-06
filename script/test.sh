#!/bin/bash
export GOPATH="$GOPATH:$(cd ..; pwd)"
cd ../src/mahjong
go test -test.bench=".*"
go test -coverprofile=covprofile
go tool cover -html=covprofile -o coverage.html
