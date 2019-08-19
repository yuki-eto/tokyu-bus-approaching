#!/usr/bin/env bash

set -e

BASE_PATH="$(cd "$(dirname $0)" && pwd)/.."

statik -src font -f
go build -v -o "${BASE_PATH}/build/tokyu_bus_approaching" "${BASE_PATH}/cmd/gui.go"
