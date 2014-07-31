package mesos

// We are reached here only from the 'Ready' state
func stateHeartbeat(d *Driver) stateFn {
	d.config.Log.Info.Println("STATE: Heartbeat")
	return stateReady
}
