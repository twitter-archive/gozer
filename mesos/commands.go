package mesos

type userCmd struct {
	command cmdCode
	task    MesosTask
}

type cmdCode int

const (
	cmdLaunch cmdCode = iota
	cmdKill
)

func (m *MesosMaster) LaunchTask(task MesosTask) error {
	m.Lock()
	defer m.Unlock()

	cmd := &userCmd{
		command: cmdLaunch,
		task:    task,
	}
	m.userCommands = append(m.userCommands, cmd)
	m.sendCommand <- len(m.userCommands)

	return nil
}
