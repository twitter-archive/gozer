package mesos

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

// We wait until HTTP Pid endpoint is ready and healthy
func stateInit(m *MesosMaster) stateFn {
	log.Print("INIT: Starting framework:", m)

	delay := time.Second
	healthURL := fmt.Sprintf("http://%s:%d/health", m.localIp, m.localPort)

	// Start Pid endpoint
	go startServing(m)

	// Now wait for healthy endpoint
	for {
		resp, err := http.Get(healthURL)
		if err == nil && resp.StatusCode == http.StatusOK {
			// Ignore content
			break
		}

		log.Print("INIT: Timeout for URL: ", healthURL, ", error =", err)
		time.Sleep(delay)
		if delay < 2*time.Minute {
			delay = delay * 2
		}
	}

	return stateRegister
}
