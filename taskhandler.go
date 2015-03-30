package flue

import (
	"time"
)

// A task is an action executed by a module
type Task struct {
	Name   string //the task name
	Node   string // The node name
	Module string
	Args   []string
	Status int // -2: queued
	// -1: running
	// >=0 : return code
	Deps      []string // A map index task1 will wait for task2, task3 and task4 to be completed
	StartTime time.Time
	EndTime   time.Time
}

// This is the structure corresponding to the "dot-graph" of a task list
// We store the nodes in a map
// The index is the source node
// The value is an array of strings containing the destination
type TaskGraphStructure struct {
	Tasks []*Task
}

func NewTask() *Task {
	return &Task{
		"null",
		"localhost",
		"dummy",
		make([]string, 1),
		-2,
		make([]string, 0),
		time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
		time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
	}

}
func NewTaskGraphStructure() *TaskGraphStructure {
	return &TaskGraphStructure{
		make([]*Task, 0),
	}
}

func GetTask(taskName string, taskStructure *TaskGraphStructure) *Task {
	for _, task := range taskStructure.Tasks {
		if task.Name == taskName {
			return task
		}
	}
	return nil
}
