package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/twitter/gozer/mesos"
)

type TaskStore struct {
	sync.RWMutex
	tasks map[string]*Task
}

func NewTaskStore() *TaskStore {
	return &TaskStore{
		tasks: make(map[string]*Task),
	}
}

func (t *TaskStore) Add(task *Task) error {
	t.Lock()
	defer t.Unlock()

	if _, ok := t.tasks[task.Id]; ok {
		return fmt.Errorf("task Id %q already exists; addition ignored", task.Id)
	}

	task.State = TaskState_INIT
	task.mesosTask = &mesos.MesosTask{
		Id:      task.Id,
		Command: task.Command,
	}
	t.tasks[task.Id] = task
	log.Printf("TASK %q State * -> %s", task.Id, task.State)

	return nil
}

func (t *TaskStore) Remove(taskId string) error {
	t.Lock()
	defer t.Unlock()

	if _, ok := t.tasks[taskId]; !ok {
		return fmt.Errorf("task Id %q not found; removal ignored", taskId)
	}

	delete(t.tasks, taskId)
	log.Printf("TASK %q removed", taskId)

	return nil
}

func (t *TaskStore) Update(taskId string, state TaskState) error {
	t.Lock()
	defer t.Unlock()

	task, ok := t.tasks[taskId]
	if !ok {
		return fmt.Errorf("task Id %q not found, update ignored", taskId)
	}

	log.Printf("TASK %q State %s -> %s", taskId, task.State, state)
	task.State = state

	return nil
}

func (t *TaskStore) Ids() []string {
	t.RLock()
	defer t.RUnlock()

	keys := make([]string, 0)
	for key := range t.tasks {
		keys = append(keys, key)
	}

	return keys
}

func (t *TaskStore) State(taskId string) (TaskState, error) {
	t.RLock()
	defer t.RUnlock()

	task, ok := t.tasks[taskId]
	if !ok {
		return "", fmt.Errorf("task Id %q not found", taskId)
	}

	return task.State, nil
}

func (t *TaskStore) MesosTask(taskId string) (*mesos.MesosTask, error) {
	t.RLock()
	defer t.RUnlock()

	task, ok := t.tasks[taskId]
	if !ok {
		return nil, fmt.Errorf("task Id %q not found", taskId)
	}

	return task.mesosTask, nil
}
