#!/bin/bash

export CGO_ENABLED=0

go test -timeout 30m -v gitlab.torproject.org/tpo/anti-censorship/pluggable-transports/webtunnel/...