package gautomator

import (
	"fmt"
	"github.com/gonum/matrix/mat64" // Matrix
	"io"
	"log"
	"time"
)

const (
    TASKQUEUED = -3
    TASKADVERTIZED = -2
    TASKRUNNING = -1

    ORPHAN = -2
    FATHER = -1
)
// A task is an action executed by a module
type Task struct {
	Id     int `json:"id"`
	Father int
	Origin string   `json:"origin"`
	Name   string   `json:"name"` //the task name
	Node   string   `json:"node"` // The node name
	Module string   `json:"module"`
	Args   []string `json:"args"`
	Status int      `json:"status"` //-3: queued
	// -2 Advertized (infored that the dependencies are done)
	// -1: running
	// >=0 : return code
	StartTime      time.Time `json:"startTime"`
	EndTime        time.Time `json:"endTime"`
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
		ORPHAN,
		"null",
		"null",
		"null",
		"dummy",
		make([]string, 1),
		TASKQUEUED,
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
	//a, b := this.AdjacencyMatrix.Dims()
	for r := 0; r < initialRowLen+addedRowLen; r++ {
		for c := 0; c < initialColLen+addedColLen; c++ {
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

func (this *TaskGraphStructure) getTaskFromName(name string) []int {
	indexA := make([]int, 1)
	indexA[0] = -1
	for _, task := range this.Tasks {
		if task.Name == name {
			if indexA[0] == -1 {
				indexA = append(indexA[1:], task.Id)
			} else {
				indexA = append(indexA, task.Id)
			}
		}
	}
	return indexA
}

func colSum(matrix *mat64.Dense, colId int) float64 {
	row, _ := matrix.Dims()
	sum := float64(0)
	for r := 0; r < row; r++ {
		sum += matrix.At(r, colId)
	}
	return sum
}

func rowSum(matrix *mat64.Dense, rowId int) float64 {
	_, col := matrix.Dims()
	sum := float64(0)
	for c := 0; c < col; c++ {
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
	// IN this array we store the row,col on which we set 1
	backup := make(map[string][]int, 0)
	_, col := this.AdjacencyMatrix.Dims()
	for _, task := range this.Tasks {
		if colSum(this.AdjacencyMatrix, task.Id) == 0 {
			id := this.getTaskFromName(task.Origin)
			// TODO There should be only one task, otherwise display an error
			if id[0] != -1 {
				// Task is a meta task
				this.Tasks[id[0]].Module = "meta"
				this.AdjacencyMatrix.Set(id[0], task.Id, float64(1))
				backup[task.Origin] = append(backup[task.Origin], id[0], task.Id)
			}
		}
		if rowSum(this.AdjacencyMatrix, task.Id) == 0 {
			id := this.getTaskFromName(task.Origin)
			// TODO There should be only one task, otherwise display an error
			if id[0] != -1 {
				for c := 0; c < col; c++ {
					add := true
					for counter := 0; counter < len(backup[task.Origin])-1; counter += 2 {
						if backup[task.Origin][counter] == id[0] && backup[task.Origin][counter+1] == c {
							add = false
						}
					}
					if add == true && this.Tasks[c].Origin != task.Origin {
						this.AdjacencyMatrix.Set(task.Id, c, this.AdjacencyMatrix.At(task.Id, c)+this.AdjacencyMatrix.At(id[0], c))
					}
				}
			}
		}
	}
	//TODO: complete the degreematrix
	return this
}

// Duplicate the task "id"
// Returns the id of the new task and the whole structure
func (this *TaskGraphStructure) DuplicateTask(name string) []int {
	newIds := make([]int, 1)
	newIds[0] = -1
	Ids := this.getTaskFromName(name)
	for _, id := range Ids {
		if id != -1 {
			newId, _ := this.AdjacencyMatrix.Dims()
			if newIds[0] == -1 {
				newIds = append(newIds[1:], newId)
			} else {
				newIds = append(newIds, newId)
			}
			newTask := NewTask()
			newTask.Id = newId
			newTask.Name = this.Tasks[id].Name
			newTask.Origin = this.Tasks[id].Origin
			newTask.Module = this.Tasks[id].Module
			newTask.Node = this.Tasks[id].Node
			newTask.Args = this.Tasks[id].Args
			newTask.Status = this.Tasks[id].Status
			this.Tasks[newId] = newTask
			this.AdjacencyMatrix = mat64.DenseCopyOf(this.AdjacencyMatrix.Grow(1, 1))
			this.DegreeMatrix = mat64.DenseCopyOf(this.DegreeMatrix.Grow(1, 1))
			for r := 0; r < newId; r++ {
				this.AdjacencyMatrix.Set(r, newId, this.AdjacencyMatrix.At(r, id))
				this.DegreeMatrix.Set(r, newId, this.DegreeMatrix.At(r, id))
			}
			// Copy the col 'id' to col 'newId'
			for c := 0; c < newId; c++ {
				this.AdjacencyMatrix.Set(newId, c, this.AdjacencyMatrix.At(id, c))
				this.DegreeMatrix.Set(newId, c, this.DegreeMatrix.At(id, c))
			}
		}
	}
	return newIds
}

// This function print the dot file associated with the graph
func (this *TaskGraphStructure) PrintDot(w io.Writer) {
	fmt.Fprintln(w, "digraph G {")
	// Writing node definition
	for _, task := range this.Tasks {
		fmt.Fprintf(w, "\t\"%v\" [\n", task.Id)
		fmt.Fprintf(w, "\t\tid = \"%v\"\n", task.Id)
		if task.Module == "meta" {
			fmt.Fprintln(w, "\t\tshape=diamond")
			fmt.Fprintf(w, "\t\tlabel=\"%v\"", task.Name)
		} else {
			fmt.Fprintf(w, "\t\tlabel = \"<name>%v|<node>%v|<module>%v\"\n", task.Name, task.Node, task.Module)
			fmt.Fprintf(w, "\t\tshape = \"record\"\n")
		}
		fmt.Fprintf(w, "\t];\n")
	}
	row, col := this.AdjacencyMatrix.Dims()
	for r := 0; r < row; r++ {
		for c := 0; c < col; c++ {
			if this.AdjacencyMatrix.At(r, c) == 1 {
				fmt.Fprintf(w, "\t%v -> %v\n", this.Tasks[r].Id, this.Tasks[c].Id)
			}
		}
	}
	fmt.Fprintln(w, "}")
}

// Returns a tasks array of all tasks with the same origin
func (this *TaskGraphStructure) getTasksWithOrigin(origin string) []Task {
	returnTasks := make([]Task, 0)
	for _, task := range this.Tasks {
		if task.Origin == origin {
			returnTasks = append(returnTasks, *task)
		}
	}
	return returnTasks
}

// Returns an extract of taskStructure containing only the tasks listed as argument.
// All other tasks are nil, and non relative elements of the matrix are zeroed
/*
func (this *TaskGraphStructure) getSubStructure(taskList []*Task) TaskGraphStructure {
	return nil
}
*/

// Duplicate the task passed as argument, and returns the new task
func (this *TaskGraphStructure) instanciate(instance TaskInstance) []*Task {
	returnTasks := make([]*Task, 0)
	// First duplicate the tasks with same name
	for _, task := range this.Tasks {
		//log.Printf("DEBUG: task %v", task.Name)
		if task.Name == instance.Taskname {
			for _, node := range instance.Hosts {
				log.Printf("DEBUG: %v", node)
				switch {
				case task.Father == FATHER:
					// Then duplicate
					log.Printf("Duplicating %v on node %v", task.Name, node)
					row, col := this.AdjacencyMatrix.Dims()
					newId := row
					newTask := NewTask()
					newTask.Father = task.Id
					newTask.Id = newId
					newTask.Name = task.Name
					newTask.Module = instance.Module
					newTask.Origin = task.Origin
					newTask.Node = node // Set the node to the new one
					newTask.Args = instance.Args
					this.Tasks[newId] = newTask
					returnTasks = append(returnTasks, newTask)
					this.AdjacencyMatrix = mat64.DenseCopyOf(this.AdjacencyMatrix.Grow(1, 1))
					for r := 0; r < row; r++ {
						for c := 0; c < col; c++ {
							if this.Tasks[r].Origin != instance.Taskname {
								this.AdjacencyMatrix.Set(r, newId, this.AdjacencyMatrix.At(r, task.Id))
							}
							if this.Tasks[c].Origin != instance.Taskname {
								this.AdjacencyMatrix.Set(newId, c, this.AdjacencyMatrix.At(task.Id, c))
							}
						}
					}
				    case task.Father == ORPHAN:
					// Do not duplicate, simply adapt
					task.Node = node
					task.Module = instance.Module
					task.Args = instance.Args
					task.Father = FATHER
				}
				// Then duplicate the tasks with same instance.Taskname
			}
		}
	}
	return returnTasks
}

// Duplicate a taskstructure
/*
func duplicateTaskGraphStructure(taskstructure *TaskGraphStructure) *TaskGraphStructure {
	return nil
}
*/

func (this *TaskGraphStructure) InstanciateTaskStructure(taskDefinition TaskDefinition) {
	for _, instance := range taskDefinition {
		//newTasks := this.duplicateTasks(instance.Taskname, host)
		log.Printf("Instance %v %v %v", instance.Taskname, instance.Hosts, instance.Module)
		if instance.Module != "" {
			instance.Module = fmt.Sprintf("%v%v", "../examples/modules/", instance.Module)
		} else {
			instance.Module = "dummy"
		}
		this.instanciate(instance)
	}
}
