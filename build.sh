#! /bin/bash

export GOPATH=`pwd`
export PATH=$PATH:$GOPATH/bin

# Check for protoc
protoc >/dev/null 2>&1
if [ $? -ne 0 ]; then
	echo "Can not find 'protoc' binary in \$PATH" 1>&2
	exit 1
fi

set -e

go get code.google.com/p/goprotobuf/{proto,protoc-gen-go}

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
