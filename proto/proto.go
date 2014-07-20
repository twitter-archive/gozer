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
	return "51e047524cf744ee257870eb479345646c0428ff"
}

// GitTag returns the Git tag that was used to build these protobuf bindings to the Mesos API.
func GitTag() string {
	return "0.19.0"
}

// GitTime returns the time of the Git commit that was used to build these protobuf bindings
// to the Mesos API.
func GitTime() time.Time {
	ts, _ := strconv.ParseInt("1401914422", 10, 64)
	return time.Unix(ts, 0)
}
