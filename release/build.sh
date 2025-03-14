#!/bin/bash

export CGO_ENABLED=0

if [ ! -d "build" ]
  then
    mkdir build
fi

pushd build
if [ ! -d "$GOARCH-$GOOS" ]
    then
      mkdir "$GOARCH-$GOOS"
fi

pushd "$GOARCH-$GOOS"

popd
popd

go build -ldflags="-s -w" -o "build/$GOARCH-$GOOS/client" gitlab.torproject.org/tpo/anti-censorship/pluggable-transports/webtunnel/main/client
go build -ldflags="-s -w" -o "build/$GOARCH-$GOOS/server" gitlab.torproject.org/tpo/anti-censorship/pluggable-transports/webtunnel/main/server

