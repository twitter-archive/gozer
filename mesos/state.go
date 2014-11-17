package mesos

// A state function is a function that does stuff, and
// then returns the next state function to be invoked.
type stateFn func(*Driver) stateFn

// Run the state machine
func (d *Driver) Run() {
	for state := stateInit; state != nil; {
		state = state(d)
	}
	// Close channels to indicate driver state machine is done.
	close(d.Updates)
	close(d.Offers)
	close(d.events)
	close(d.command)
}

func stateStop(d *Driver) stateFn {
	d.config.Log.Info.Println("STOP: Stopping framework:", d)
	return nil
}
