#!/bin/bash
export GOPATH="$GOPATH:$(cd ..; pwd)"
cd ../src/mahjong
go test -test.bench=".*"