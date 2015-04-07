package flue

import (
	"fmt"
	"github.com/gonum/matrix/mat64" // Matrix
	"log"
	"time"
)

// A task is an action executed by a module
type Task struct {
	Id     int
	Origin string
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
	// merging adjacency matrix
	initialRowLen, initialColLen := this.AdjacencyMatrix.Dims()
	addedRowLen, addedColLen := taskStructure.AdjacencyMatrix.Dims()
	this.AdjacencyMatrix = mat64.DenseCopyOf(this.AdjacencyMatrix.Grow(addedRowLen, addedColLen))
	a,b:=this.AdjacencyMatrix.Dims()
	log.Printf("DEBUG: %v,%v",a,b)
	for r := 0; r < initialRowLen+addedRowLen; r++ {
		for c := 0; c < initialColLen+addedColLen; c++ {
			log.Printf("Adjacency r:%v,c:%v", r, c)
			switch {
			case r < initialRowLen && c < initialColLen:
				// If we are in the original matrix: do nothing
			case r < initialRowLen && c > initialColLen:
				// If outside, put some zero
				this.AdjacencyMatrix.Set(r, c, float64(0))
			case r > initialRowLen && c < initialColLen:
				// If outside, put some zero
				this.AdjacencyMatrix.Set(r, c, float64(0))
			case r >= initialRowLen && c >= initialColLen:
				// Add the new matrix
				this.AdjacencyMatrix.Set(r, c, taskStructure.AdjacencyMatrix.At(r-initialRowLen, c-initialColLen))
			}
		}
	}
	// merging degree matrix
	initialRowLen, initialColLen = this.DegreeMatrix.Dims()
	addedRowLen, addedColLen = taskStructure.DegreeMatrix.Dims()
	this.DegreeMatrix = mat64.DenseCopyOf(this.DegreeMatrix.Grow(addedRowLen, addedColLen))
	for r := 0; r < initialRowLen+addedRowLen; r++ {
		for c := 0; c < initialColLen+addedColLen; c++ {
			log.Printf("Degree r:%v,c:%v", r, c)
			switch {
			case r < initialRowLen && c < initialColLen:
				// If we are in the original matrix: do nothing
			case r < initialRowLen && c > initialColLen:
				// If outside, set zero
				this.DegreeMatrix.Set(r, c, float64(0))
			case r > initialRowLen && c < initialColLen:
				// If outside, set zero
				this.DegreeMatrix.Set(r, c, float64(0))
			case r >= initialRowLen && c >= initialColLen:
				// Add the new matrix
				this.DegreeMatrix.Set(r, c, taskStructure.DegreeMatrix.At(r-initialRowLen, c-initialColLen))
			}
		}
	}
	actualSize := len(this.Tasks)
	for i, task := range taskStructure.Tasks {
		task.Id = actualSize + i
		this.Tasks[actualSize+i] = task
	}
	return this
}

func (this *TaskGraphStructure) getTaskFromName(name string) (int, *Task) {
    for index, task := range this.Tasks {
	if task.Name == name {
	    return index, task
	}
    }
    return -1, nil
}

func colSum(matrix *mat64.Dense, colId int) float64 {
    row, _ := matrix.Dims()
    sum := float64(0)
    for r := 0 ; r < row ; r++ {
	sum += matrix.At(r, colId)
    }
    return sum
}

func rowSum(matrix *mat64.Dense, rowId int) float64 {
    _, col := matrix.Dims()
    sum := float64(0)
    for c := 0 ; c < col ; c++ {
	sum += matrix.At(rowId, c)	
    }
    return sum
}
// the aim of this function is to find if a task has a subdefinition (aka an origin) and change it
// Example:
// imagine the graphs
// digraph bla {
//    a -> b;
//    b -> c;
// }
// digraph b {
//    alpha -> gamma;
// }
// then alpha and beta will have "b" as Origin.
// therefore we should add a link in the AdjacencyMatix and in the DegreeMatrix
func (this *TaskGraphStructure) Relink() *TaskGraphStructure {
    _, col := this.AdjacencyMatrix.Dims()
    for _, task := range this.Tasks {
	if colSum(this.AdjacencyMatrix, task.Id) == 0 {
	    id, _ := this.getTaskFromName(task.Origin)
	    if id != -1 {
		this.AdjacencyMatrix.Set(id, task.Id, float64(1))
	    }
	}
	if rowSum(this.AdjacencyMatrix, task.Id) == 0 {
	    id, _ := this.getTaskFromName(task.Origin)
	    if id != -1 {
		for c :=0 ; c < col ; c++ {
		    this.AdjacencyMatrix.Set(task.Id,c,this.AdjacencyMatrix.At(task.Id,c) + this.AdjacencyMatrix.At(id,c))
		}
	    }
	}
    }
    /*
    for each task {
          // Relink the origin 
          if sum(col(task)) = 0 {
	       set 1 at Adjacency(rowId(origin),colId(task))
	   }
	   // Relink the destination
	   if sum(row(task)) = 0 {
		row(task)=row(task)+row(origin)
	   }
       }

    */
    return this
}
