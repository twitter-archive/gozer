Protobuf compilation
====================

The Mesos core uses Google's Protobuf to serialize and de-serialize messages between various portions of the software.  The infrastructure here exists to collect and create a set of `proto` files that are then compiled into a set of golang source packages.  In order to make using these package easier, the output from this directory is then commited and added to http://github.com/twitter/gozer/proto in order for other golang sources to use.

Prerequisites & Setup
---------------------

The build scripts within this directory require access to a number of protobuf related compilers.  In particular at least version 2.5.0 of [Protocol Buffers](https://code.google.com/p/protobuf/) will be required, as well as [goprotobuf](https://code.google.com/p/goprotobuf/) support.

The easiest setup for an OSX environment will likely be to use [Homebrew](http://brew.sh/) to install the `protobuf` support.  After that, if you have a fully configured golang environment, `go get` can be used to install the required golang support:

```bash
shell$ brew install protobuf
shell$ cd $GOPATH
shell$ go get code.google.com/p/goprotobuf/proto
shell$ go get code.google.com/p/goprotobuf/protoc-gen-go
shell$ go install code.google.com/p/goprotobuf/protoc-gen-go
shell$ ls -al $GOPATH/bin/protoc-gen-go
```

Building Protobuf definitions
-----------------------------

TODO(weingart): more information here
