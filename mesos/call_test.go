package mesos

import "testing"

import "github.com/twitter/gozer/proto/scheduler.pb"

func TestPath(t *testing.T) {
	in := mesos_scheduler.Call_REGISTER
	const out = "mesos.internal.RegisterFrameworkMessage"

	call := &mesos_scheduler.Call{
		Type: &in,
	}

	got, err := path(call)
	if err != nil || got != out {
		t.Errorf("path(%v): got %v, want %v", in, got, out)
	}
}
