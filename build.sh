#!/bin/bash

export GOPATH=`pwd`
export PATH=$PATH:$GOPATH/bin

set -e

for protobuf in mesos messages scheduler; do
	OUT_PATH=$GOPATH/src/$protobuf.pb

	mkdir -p $OUT_PATH

	cd $GOPATH/proto
	protoc --go_out=$OUT_PATH/ $protobuf.proto

	cd $GOPATH
	go install $protobuf.pb
done

cd $GOPATH
go install mesos
go build gozer
