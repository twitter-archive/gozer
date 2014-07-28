package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/twitter/gozer/mesos"
	mesos_pb "github.com/twitter/gozer/proto/mesos.pb"
)

var (
	user       = flag.String("user", "", "The user to register as")
	port       = flag.Int("port", 4343, "Port to listen on for the API endpoint")
	master     = flag.String("master", "localhost", "Hostname of the master")
	masterPort = flag.Int("masterPort", 5050, "Port of the master")

	taskstore = NewTaskStore()

	// TODO(dhamon): flags for log level
	log = mesos.NewLog(mesos.LogConfig{
		Prefix: "gozer",
		Info:   os.Stdout,
		Warn:   os.Stderr,
		Error:  os.Stderr},
	)
)

type TaskState string

const (
	TaskState_INIT     TaskState = "INIT"
	TaskState_STARTING TaskState = "STARTING"
	TaskState_RUNNING  TaskState = "RUNNING"
	TaskState_FINISHED TaskState = "FINISHED"
	TaskState_FAILED   TaskState = "FAILED"
	TaskState_KILLED   TaskState = "KILLED"
	TaskState_LOST     TaskState = "LOST"
)

var taskStateMap = map[mesos_pb.TaskState]TaskState{
	mesos_pb.TaskState_TASK_STAGING:  TaskState_STARTING,
	mesos_pb.TaskState_TASK_STARTING: TaskState_STARTING,
	mesos_pb.TaskState_TASK_RUNNING:  TaskState_RUNNING,
	mesos_pb.TaskState_TASK_FINISHED: TaskState_FINISHED,
	mesos_pb.TaskState_TASK_FAILED:   TaskState_FAILED,
	mesos_pb.TaskState_TASK_KILLED:   TaskState_KILLED,
	mesos_pb.TaskState_TASK_LOST:     TaskState_LOST,
}

type Task struct {
	Id        string           `json:"id"`
	Command   string           `json:"command"`
	State     TaskState        `json:"state"`
	mesosTask *mesos.MesosTask `json:"-"`
	// TODO(dhamon): resource requirements
}

func main() {
	flag.Parse()

	go startHTTP()

	log.Info.Println("Registering")
	driver, err := mesos.New("gozer", *user, *master, *masterPort)
	if err != nil {
		log.Error.Fatal(err)
	}

	// Shepherd all our tasks
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
	// For now we use a simple loop to do a very naive management of tasks, updates, events,
	// errors, etc.
	for {

		select {

		case update := <-driver.Updates:
			log.Info.Println("Status update:", update)
			state, err := taskstore.State(update.TaskId)
			if err != nil {
				log.Error.Printf("Failed to get current state for updated task %q: %+s", update.TaskId, err)
				continue
			}

			newState, ok := taskStateMap[update.State]
			if !ok {
				log.Error.Printf("Unknown mesos task state: %s", update.State)
				continue
			}

			log.Info.Println("Updating task state from", state, "to", newState)

			switch newState {
			case TaskState_FAILED:
			case TaskState_FINISHED:
			case TaskState_KILLED:
			case TaskState_LOST:
				taskstore.Remove(update.TaskId)

			default:
				taskstore.Update(update.TaskId, newState)
			}

			update.Ack()

		case <-time.After(5 * time.Second):
			log.Info.Println("Gozer: checking for tasks")
			// After a timeout, see if there any tasks to launch
			taskIds := taskstore.Ids()
			for _, taskId := range taskIds {
				state, err := taskstore.State(taskId)
				if err != nil {
					log.Error.Printf("Error getting task state for task %q: %+v", taskId, err)
					continue
				}

				if state != TaskState_INIT {
					continue
				}

				mesosTask, err := taskstore.MesosTask(taskId)
				if err != nil {
					log.Error.Printf("Error getting mesos task for task %q: %+v", taskId, err)
					continue
				}

				// Start this task (very naive method)
				offer := <-driver.Offers

				// TODO(dhamon): Check for resources before launching
				err = driver.LaunchTask(offer, mesosTask)
				if err != nil {
					log.Error.Printf("Error launching task %q: %+v", taskId, err)
					continue
				}

				taskstore.Update(taskId, TaskState_STARTING)
			}
		}
	}
}
