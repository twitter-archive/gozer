package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"mesos"
)

const (
	frameworkName = "gozer"
)

var (
	user = flag.String("user", "", "The user to register as")
	port = flag.Int("port", 4343, "Port to listen on for HTTP endpoint")

	tasks = make(chan Task)
)

type Task struct {
	Command string
	// TODO(dhamon): resource requirements
}

func startAPI() {
	log.Printf("api listening on port %d", *port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil); err != nil {
		log.Fatalf("failed to start listening on port %d", *port)
	}
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.Header().Add("Allow", "POST")
		w.WriteHeader(http.StatusMethodNotAllowed)
		log.Printf("received addtask request with unexpected method. want %q, got %q: %+v", "POST", r.Method, r)
	}

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("ERROR: failed to read body from addtask request %+v: %+v", r, err)
		// TODO(dhamon): Better error for this case.
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var task Task
	err = json.Unmarshal(body, &task)

	if err != nil {
		log.Printf("ERROR: failed to parse JSON body from addtask request %+v: %+v", r, err)
		// TODO(dhamon): Better error for this case.
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	tasks <- task

	w.WriteHeader(http.StatusOK)
}

func main() {
	flag.Parse()

	http.HandleFunc("/api/addtask", addTaskHandler)
	go startAPI()

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	log.Printf("Registering...")
	err := mesos.Register(*user, frameworkName)
	if err != nil {
		log.Fatal(err)
	}

	taskId := 0
	for {
		// TODO(dhamon): wait for offers in go routine
		log.Printf("Waiting for offers...")
		offers, err := mesos.WaitForOffers()
		if err != nil {
			log.Fatal(err)
		}

		// TODO(dhamon): add offers to list and consume as tasks come in.
		for _, offer := range offers {
			log.Printf("Received offer %+v", offer)
			log.Printf("Waiting for tasks...")
			task := <-tasks

			// TODO(dhamon): Decline offers if resources don't match.
			err = mesos.LaunchTask(offer, fmt.Sprintf("gozer-task-%d", taskId), task.Command)
			taskId += 1
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
