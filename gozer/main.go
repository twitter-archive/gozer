package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/twitter/gozer/mesos"
)

var (
	user       = flag.String("user", "", "The user to register as")
	port       = flag.Int("port", 4343, "Port to listen on for HTTP endpoint")
	master     = flag.String("master", "localhost", "Hostname of the master")
	masterPort = flag.Int("masterPort", 5050, "Port of the master")

	taskstore = &TaskStore{
		Tasks: make(map[string]*Task),
	}
)

type TaskState int

const (
	TaskState_UNKNOWN TaskState = iota
	TaskState_INIT
	TaskState_STARTING
	TaskState_RUNNING
	TaskState_FAILED
	TaskState_FINISHED
)

func (t *TaskState) String() string {
	switch *t {
	case TaskState_INIT:
		return "INIT"
	case TaskState_STARTING:
		return "STARTING"
	case TaskState_RUNNING:
		return "RUNNING"
	case TaskState_FAILED:
		return "FAILED"
	case TaskState_FINISHED:
		return "FINISHED"
	default:
		return "UNKNOWN"
	}
}

type Task struct {
	Id        string           `json:"id"`
	Command   string           `json:"command"`
	State     TaskState        `json:"-"`
	MesosTask *mesos.MesosTask `json:"-"`
	// TODO(dhamon): resource requirements
}

type TaskStore struct {
	sync.RWMutex
	Tasks map[string]*Task
}

func (t *TaskStore) AddTask(task *Task) error {
	t.Lock()
	defer t.Unlock()

	if _, ok := t.Tasks[task.Id]; ok {
		return fmt.Errorf("Task Id '%s' already exists, addition ignored", task.Id)
	}

	task.State = TaskState_INIT
	task.MesosTask = &mesos.MesosTask{
		Id:      task.Id,
		Command: task.Command,
	}
	t.Tasks[task.Id] = task
	log.Print("TASK '%s' State * -> %s", task.Id, task.State)

	return nil
}

func (t *TaskStore) UpdateTask(taskId string, state TaskState) error {
	t.Lock()
	defer t.Unlock()

	task, ok := t.Tasks[taskId]
	if !ok {
		return fmt.Errorf("Task Id '%s' not found, update ignored", taskId)
	}

	log.Printf("Task '%s' State %s -> %s", taskId, &task.State, &state)
	task.State = state

	return nil
}

func (t *TaskStore) TaskIds() []string {
	t.RLock()
	defer t.RUnlock()

	keys := make([]string, 0)
	for key := range t.Tasks {
		keys = append(keys, key)
	}

	return keys
}

func (t *TaskStore) TaskState(taskId string) (TaskState, error) {
	t.RLock()
	defer t.RUnlock()

	task, ok := t.Tasks[taskId]
	if !ok {
		return TaskState_UNKNOWN, fmt.Errorf("Task Id '%s' not found", taskId)
	}

	return task.State, nil
}

func (t *TaskStore) MesosTask(taskId string) (*mesos.MesosTask, error) {
	t.RLock()
	defer t.RUnlock()

	task, ok := t.Tasks[taskId]
	if !ok {
		return nil, fmt.Errorf("Task Id '%s' not found", taskId)
	}

	return task.MesosTask, nil
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

	taskstore.AddTask(&task)

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

	// Launch simple /bin/true after a little while
	go func() {
		time.Sleep(10 * time.Second)

		trueTask := &Task{
			Id:      "gozer-bin-true",
			Command: "/bin/true",
		}

		taskstore.AddTask(trueTask)
	}()

	// Shephard all our tasks
	//
	// Note: This will require a significant re-architecting, most likely to break out the
	// gozer based tasks and their state transitions, possibly using a go-routine per task,
	// which may limit the total number of tasks we can handle (100k go-routines might be
	// too much).  It should make for a simple abstraction, where the update routine should
	// then be able to simply use a channel to post a state transition to the gozer task,
	// which the gozer task manager (per task) go routine would then use to transition the
	// task through its state diagram.  This should also make it very simple to detect bad
	// transitions.
	//
	// It would be nice if we could just only use the mesos.TaskState_TASK_* states, however,
	// they do not encompass the ideas of PENDING (waiting for offers), and ASSIGNED (offer
	// selected, waiting for running), nor do they encompass tear-down and death.
	//
	// For now we use a simple loop to do a very naiive management of tasks, updates, events,
	// errors, etc.
	for {

		select {
		// case <-master.Events:
		// case <-master.Updates:

		case <-time.After(5 * time.Second):
			log.Print("Gozer: timeout, check tasks")
			// After a timeout, see if there any tasks to launch
			taskIds := taskstore.TaskIds()
			for _, taskId := range taskIds {
				state, err := taskstore.TaskState(taskId)
				if err != nil || state != TaskState_INIT {
					continue
				}

				mesosTask, err := taskstore.MesosTask(taskId)
				if err != nil {
					log.Print(err)
					continue
				}

				// Start this task (very naiive method)
				offer := <-master.Offers
				err = master.LaunchTask(offer, mesosTask)
				if err != nil {
					log.Print(err)
					continue
				}

				taskstore.UpdateTask(taskId, TaskState_STARTING)
			}
		}
	}
}
