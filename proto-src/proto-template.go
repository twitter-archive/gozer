// DO NOT EDIT!
// This is file is automatically generated!

package proto

import (
	"time"
)

func GitSHA() string {
	return "@GIT_SHA@"
}

func GitTime() time.Time {
	return time.Unix(@GIT_TS@, 0)
}
