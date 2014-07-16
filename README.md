gozer
=====

Prototype mesos framework using new low-level API built in Go.

Dependencies
------------

This software requires a number of dependencies in order for the build to work. In particular, the [Google Protobuf](https://developers.google.com/protocol-buffers) infrastructure needs to be present.  On OSX, the easiest way to likely get working version of `protoc` is by using `brew install protobuf`.

Building
--------

Once the required dependencies are installed and available on your `$PATH`, the quickest way to build this is to:

```bash
shell$ git clone https://github.com/dominichamon/gozer
shell$ cd gozer
gozer$ ./build.sh
gozer$ ./test.sh
gozer$ bin/gozer --help
```
