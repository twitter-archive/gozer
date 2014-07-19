// This is here to have 'go test ...' "ignore" this directory

package proto

import (
	"testing"
)

func TestIgnored(t *testing.T) {
	t.Log("This test is here to be ignored.")
}
