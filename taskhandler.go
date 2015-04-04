package flue

import (
	"fmt"
	"github.com/gonum/matrix/mat64" // Matrix
	"time"
)

// A task is an action executed by a module
type Task struct {
	Id     int
	Name   string //the task name
	Node   string // The node name
	Module string
	Args   []string
	Status int // -2: queued
	// -1: running
	// >=0 : return code
	StartTime      time.Time
	EndTime        time.Time
	TaskCanRunChan chan bool // true: run, false: wait
}

// This is the structure corresponding to the "dot-graph" of a task list
// We store the nodes in a map
// The index is the source node
type TaskGraphStructure struct {
	Tasks           map[int]*Task
	DegreeMatrix    *mat64.Dense
	AdjacencyMatrix *mat64.Dense // Row id is the map id of the source task
	// Col id is the map id of the destination task
}

func (this *TaskGraphStructure) PrintAdjacencyMatrix() {
	rowSize, colSize := this.AdjacencyMatrix.Dims()
	fmt.Printf("  ")
	for c := 0; c < colSize; c++ {
		fmt.Printf("%v ", this.Tasks[c].Name)
	}
	fmt.Printf("\n")
	for r := 0; r < rowSize; r++ {
		fmt.Printf("%v ", this.Tasks[r].Name)
		for c := 0; c < colSize; c++ {
			fmt.Printf("%v ", this.AdjacencyMatrix.At(r, c))
		}
		fmt.Printf("\n")
	}
}

func (this *TaskGraphStructure) PrintDegreeMatrix() {
	rowSize, colSize := this.DegreeMatrix.Dims()
	for r := 0; r < rowSize; r++ {
		for c := 0; c < colSize; c++ {
			fmt.Printf("%v ", this.DegreeMatrix.At(r, c))
		}
		fmt.Printf("\n")
	}
}

func NewTask() *Task {
	return &Task{
		-1,
		"null",
		"localhost",
		"dummy",
		make([]string, 1),
		-2,
		time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
		time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
		make(chan bool),
	}

}
func NewTaskGraphStructure() *TaskGraphStructure {
	return &TaskGraphStructure{
		make(map[int]*Task, 0),
		mat64.NewDense(0, 0, nil),
		mat64.NewDense(0, 0, nil),
	}
}

// Returns a combination of the current structure
// and the one passed as argument
func (this *TaskGraphStructure) AugmentTaskStructure(taskStructure *TaskGraphStructure) *TaskGraphStructure {

}
