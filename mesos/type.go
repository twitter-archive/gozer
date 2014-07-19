package mesos

import (
	"sync"
)

type MesosMaster struct {
	sync.RWMutex

	sendCommand   chan int
	receivedEvent chan int

	userCommands []*userCmd
}

type MesosTask struct {
	Command string
}
