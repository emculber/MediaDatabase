package main

import (
	"errors"
	"time"
)

type Task struct {
	Id        int
	Name      string
	Completed bool
	Due       time.Time
	ParentId  int
}

type TaskTree struct {
	Task     Task
	Children []TaskTree
}

func (task *Task) OK() error {
	if len(task.Name) == 0 {
		return errors.New("No Task Name")
	}
	return nil
}
