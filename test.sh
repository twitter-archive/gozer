#! /bin/bash

export GOPATH=`pwd`
export PATH=$PATH:$GOPATH/bin

set -e

go test mesos
go test gozer

