package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/twitter/gozer/mesos"
)

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
	log.Printf("TASK '%s' State * -> %s", task.Id, task.State)

	return nil
}

func (t *TaskStore) UpdateTask(taskId string, state TaskState) error {
	t.Lock()
	defer t.Unlock()

	task, ok := t.Tasks[taskId]
	if !ok {
		return fmt.Errorf("Task Id '%s' not found, update ignored", taskId)
	}

	log.Printf("Task '%s' State %s -> %s", taskId, task.State, state)
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
