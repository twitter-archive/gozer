package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

var taskIndex = 0

func startHTTP() {
	http.HandleFunc("/api/addtask", addTaskHandler)
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

	var task Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		log.Printf("ERROR: failed to parse JSON body from addtask request %+v: %+v", r, err)
		// TODO(dhamon): Better error for this case.
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(task.Id) == 0 {
		task.Id = fmt.Sprintf("gozer-task-%d", taskIndex)
		taskIndex += 1
	}

	taskstore.Add(&task)

	w.WriteHeader(http.StatusOK)
}
