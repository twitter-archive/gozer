package mesos

// We are reached here only from the 'Ready' state
func stateHeartbeat(d *Driver) stateFn {
	d.log.Info.Println("STATE: Heartbeat")
	return stateReady
}
