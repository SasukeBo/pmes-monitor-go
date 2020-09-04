#!/bin/bash

echo "change work directory to" $GOPATH/src/github.com/SasukeBo/pmes-device-monitor ...
cd $GOPATH/src/github.com/SasukeBo/pmes-device-monitor

echo "start service ..."
go run server.go
