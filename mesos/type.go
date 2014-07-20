package mesos

import (
	"sync"

	"github.com/twitter/gozer/proto/mesos.pb"
)

type MesosMasterLocation struct {
	Hostname string
	Port     int
}

type MesosMasterConfig struct {
	FrameworkName  string
	RegisteredUser string
	Masters        []MesosMasterLocation
}

type MesosMaster struct {
	sync.RWMutex

	config    MesosMasterConfig
	localIp   string
	localPort int

	frameworkId mesos.FrameworkID

	sendCommand   chan int
	receivedEvent chan int

	userCommands []*userCmd
}

type MesosTask struct {
	Command string
}
