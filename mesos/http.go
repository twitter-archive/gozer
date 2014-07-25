package mesos

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func startServing(d *Driver) {

	// TODO(weingart): Grab an emphemeral port for this instead and toss it into MesosMaster
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(rw http.ResponseWriter, req *http.Request) {
		fmt.Fprint(rw, "OK\r\n")
	})
	mux.Handle("/", d)

	log.Printf("listening on port %d", d.localPort)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", d.localPort), mux); err != nil {
		log.Fatalf("failed to start listening on port %d", d.localPort)
	}
}

func (d *Driver) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.Header().Add("Allow", "POST")
		w.WriteHeader(http.StatusMethodNotAllowed)
		log.Printf("received request with unexpected method. want %q, got %q: %+v", "POST", r.Method, r)
		return
	}

	pathElements := strings.Split(r.URL.Path, "/")

	if pathElements[1] != d.config.FrameworkName {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(fmt.Sprintf("unexpected path. want %q, got %q", d.config.FrameworkName, pathElements[1])))
		log.Printf("received request with unexpected path. want %q, got %q: %+v", d.config.FrameworkName, pathElements[1], r)
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

	d.events <- event

	w.WriteHeader(http.StatusOK)
}
