// DO NOT EDIT!  This is file is automatically generated!

/*
Package proto is automatically generated and is the lowest level interface to the Mesos system.
There are no user-servicable parts in here.

The only semi-public functions are here to introspect the generated version information.
*/
package proto

import (
	"strconv"
	"time"
)

// GitSHA returns the Git SHA that was used to build these protobuf bindings to the Mesos API.
func GitSHA() string {
	return "fe0a39112f3304283f970f1b08b322b1e970829d"
}

// GitTag returns the Git tag that was used to build these protobuf bindings to the Mesos API.
func GitTag() string {
	return "0.20.1"
}

// GitTime returns the time of the Git commit that was used to build these protobuf bindings
// to the Mesos API.
func GitTime() time.Time {
	ts, _ := strconv.ParseInt("1411074811", 10, 64)
	return time.Unix(ts, 0)
}
