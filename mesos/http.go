package mesos

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

const (
	port = 8888
)

// TODO(weingart): create a htt.ServeMux and register a healthcheck URI on this server,
// which can then be used by the state machine to wait until this endpoint is up and ready
// to receive events/calls from the master.  The original SyncGroup was racy too.
func startServing(m *MesosMaster) {
	log.Printf("listening on port %d", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), m); err != nil {
		log.Fatalf("failed to start listening on port %d", port)
	}
}

func (m *MesosMaster) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.Header().Add("Allow", "POST")
		w.WriteHeader(http.StatusMethodNotAllowed)
		log.Printf("received request with unexpected method. want %q, got %q: %+v", "POST", r.Method, r)
		return
	}

	pathElements := strings.Split(r.URL.Path, "/")

	if pathElements[1] != m.config.FrameworkName {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(fmt.Sprintf("unexpected path. want %q, got %q", m.config.FrameworkName, pathElements[1])))
		log.Printf("received request with unexpected path. want %q, got %q: %+v", m.config.FrameworkName, pathElements[1], r)
		return
	}

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("ERROR: failed to read body from request %+v: %+v", r, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	event, err := bytesToEvent(pathElements[2], body)
	if err != nil {
		log.Printf("ERROR: %+v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	events <- event

	w.WriteHeader(http.StatusOK)
}
