// DO NOT EDIT!
// This is file is automatically generated!

package proto

import (
	"strconv"
	"time"
)

func GitSHA() string {
	return "@GIT_SHA@"
}

func GitTime() time.Time {
	ts, _ := strconv.ParseInt("@GIT_TS@", 10, 64)
	return time.Unix(ts, 0)
}
