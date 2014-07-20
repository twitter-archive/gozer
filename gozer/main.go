package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/twitter/gozer/mesos"
)

var (
	user       = flag.String("user", "", "The user to register as")
	port       = flag.Int("port", 4343, "Port to listen on for HTTP endpoint")
	master     = flag.String("master", "localhost", "Hostname of the master")
	masterPort = flag.Int("masterPort", 5050, "Port of the master")

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
	master, err := mesos.New(&mesos.MesosMasterConfig{
		FrameworkName:  "gozer",
		RegisteredUser: "nobody",
		Masters: []mesos.MesosMasterLocation{mesos.MesosMasterLocation{
			Hostname: *master,
			Port:     *masterPort,
		}},
	})
	if err != nil {
		log.Fatal(err)
	}

	// Start Framework engine
	go master.Run()

	// TODO(dhamon): add event loop to respond to tasks by launching if offers available.
	taskId := 0
	for {

		// TODO(dhamon): add offers to list and consume as tasks come in.
		for offer := range master.Offers {
			log.Printf("Received offer %+v", *offer)

			// Launch debug "/bin/true" task
			if taskId == 0 {
				trueTask := &mesos.MesosTask{
					Id:      fmt.Sprintf("gozer-task-%d", taskId),
					Command: "/bin/true",
				}
				err = master.LaunchTask(*offer, trueTask)
				if err != nil {
					log.Fatal(err)
				}
				continue
			}

			log.Printf("Waiting for tasks...")
			task := <-tasks

			// TODO(dhamon): Decline offers if resources don't match.
			err = master.LaunchTaskOld(*offer, fmt.Sprintf("gozer-task-%d", taskId), task.Command)
			taskId += 1
			if err != nil {
				log.Fatal(err)
			}

		}
	}
}
