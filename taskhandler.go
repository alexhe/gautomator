package gautomator

import (
	"fmt"
	"github.com/gonum/matrix/mat64" // Matrix
	"io"
	"log"
	"time"
)

const (
	TASKQUEUED     = -3
	TASKADVERTIZED = -2
	TASKRUNNING    = -1

	ORPHAN = -2
	FATHER = -1
)

// A task is an action executed by a module
type Task struct {
	Id       int `json:"id"`
	Father   int
	OriginId int      `json:"originId"`
	Origin   string   `json:"origin"`
	Name     string   `json:"name"` //the task name
	Node     string   `json:"node"` // The node name
	Module   string   `json:"module"`
	Args     []string `json:"args"`
	Status   int      `json:"status"` //-3: queued
	// -2 Advertized (infored that the dependencies are done)
	// -1: running
	// >=0 : return code
	StartTime      time.Time `json:"startTime"`
	EndTime        time.Time `json:"endTime"`
	TaskCanRunChan chan bool // true: run, false: wait
	Debug          string
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
		-1,
		"null",
		"null",
		"null",
		"dummy",
		make([]string, 1),
		TASKQUEUED,
		time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
		time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
		make(chan bool),
		"null",
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
		id := this.getTaskFromName(task.Origin)
		if id[0] != -1 && task.OriginId == -1 {
			task.OriginId = id[0]
		}
		if colSum(this.AdjacencyMatrix, task.Id) == 0 {
			// TODO There should be only one task, otherwise display an error
			if task.OriginId != -1 {
				// Task is a meta task
				this.Tasks[task.OriginId].Module = "meta"
				this.AdjacencyMatrix.Set(task.OriginId, task.Id, float64(1))
				backup[task.Origin] = append(backup[task.Origin], task.OriginId, task.Id)
			}
		}
		if rowSum(this.AdjacencyMatrix, task.Id) == 0 {
			// TODO There should be only one task, otherwise display an error
			if task.OriginId != -1 {
				for c := 0; c < col; c++ {
					add := true
					for counter := 0; counter < len(backup[task.Origin])-1; counter += 2 {
						if backup[task.Origin][counter] == task.OriginId && backup[task.Origin][counter+1] == c {
							add = false
						}
					}
					if add == true && this.Tasks[c].Origin != task.Origin {
						this.AdjacencyMatrix.Set(task.Id, c, this.AdjacencyMatrix.At(task.Id, c)+this.AdjacencyMatrix.At(task.OriginId, c))
					}
				}
			}
		}
	}
	//TODO: complete the degreematrix
	return this
}

// This function print the dot file associated with the graph
func (this *TaskGraphStructure) PrintDot(w io.Writer) {
	fmt.Fprintln(w, "digraph G {")
	// Writing node definition
	for _, task := range this.Tasks {
		fmt.Fprintf(w, "\t\"%v\" [\n", task.Id)
		fmt.Fprintf(w, "\t\tid = \"%v\"\n", task.Id)
		//		if task.Module == "meta" {
		//			fmt.Fprintln(w, "\t\tshape=diamond")
		//			fmt.Fprintf(w, "\t\tlabel=\"%v\"", task.Name)
		//		} else {
		fmt.Fprintf(w, "\t\tlabel = \"<name>%v(%v/%v/%v)|<node>%v|<module>%v|<debug>%v\"\n", task.Name, task.Id, task.OriginId, task.Origin, task.Node, task.Module, task.Debug)
		fmt.Fprintf(w, "\t\tshape = \"record\"\n")
		//		}
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

// Duplicate the task passed as argument, and returns the new task
func (this *TaskGraphStructure) instanciate(instance TaskInstance) []*Task {
	returnTasks := make([]*Task, 0)
	// First duplicate the tasks with same name
	for _, task := range this.Tasks {
		if task.Name == instance.Taskname {
			for _, node := range instance.Hosts {
				switch {
				case task.Father == FATHER:
					// Then duplicate
					log.Printf("Duplicating %v on node %v", task.Name, node)
					row, col := this.AdjacencyMatrix.Dims()
					newId := row
					newTask := NewTask()
					newTask.Father = task.Id
					newTask.OriginId = task.Id
					newTask.Id = newId
					newTask.Name = task.Name
					if task.Module != "meta" {
						newTask.Module = instance.Module
					}
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
					this = this.AugmentTaskStructure(this.duplicateSubtasks(newTask, node, instance))
				case task.Father == ORPHAN:
					// Do not duplicate, simply adapt
					task.Node = node
					if task.Module == "meta" {
						task.Module = instance.Module
					}
					task.Args = instance.Args
					this.adaptSubtask(task, node, instance)
					task.Father = FATHER
				}
				// Then duplicate the tasks with same Father
			}
		}
	}
	return returnTasks
}

func (this *TaskGraphStructure) adaptSubtask(father *Task, node string, instance TaskInstance) {
	_, col := this.AdjacencyMatrix.Dims()
	for c := 0; c < col; c++ {
		if this.AdjacencyMatrix.At(father.Id, c) == 1 && this.Tasks[c].OriginId == father.Id {
			this.Tasks[c].Father = father.Id
			this.Tasks[c].Node = node
			this.Tasks[c].Debug = "adaptSubtask"
			this.Tasks[c].Module = instance.Module
			this.Tasks[c].Args = instance.Args
		}
	}

}
func (this *TaskGraphStructure) duplicateSubtasks(father *Task, node string, instance TaskInstance) *TaskGraphStructure {
	taskStructure := NewTaskGraphStructure()
	// Get all the tasks with origin father.Name
	myindex := 0
	// Define a new origin composed of the Id
	for _, task := range this.Tasks {
		if father.Father == task.OriginId {
			// task match, create a new task with the same informations...
			newTask := NewTask()
			newTask.Id = myindex
			newTask.Name = task.Name
			newTask.Node = node
			newTask.Father = task.Id
			newTask.OriginId = father.Id
			newTask.Origin = father.Name
			newTask.Module = instance.Module
			newTask.Args = instance.Args
			newTask.Debug = "duplicateSubtasks"
			// ... Add it to the structure...
			taskStructure.Tasks[myindex] = newTask
			// ... Extract the matrix associated
			taskStructure.AdjacencyMatrix = mat64.DenseCopyOf(taskStructure.AdjacencyMatrix.Grow(1, 1))
			taskStructure.DegreeMatrix = mat64.DenseCopyOf(taskStructure.DegreeMatrix.Grow(1, 1))
			myindex += 1
			// And add it to the structure as well
		}
	}
	row, col := taskStructure.AdjacencyMatrix.Dims()
	for r := 0; r < row; r++ {
		for c := 0; c < col; c++ {
			taskStructure.AdjacencyMatrix.Set(r, c, this.AdjacencyMatrix.At(taskStructure.Tasks[r].Father, taskStructure.Tasks[c].Father))
		}
	}
	return taskStructure
}

func (this *TaskGraphStructure) InstanciateTaskStructure(taskDefinition TaskDefinition) {
	for _, instance := range taskDefinition {
		//newTasks := this.duplicateTasks(instance.Taskname, host)
		log.Printf("Instance %v %v %v", instance.Taskname, instance.Hosts, instance.Module)
		if instance.Module != "" {
			instance.Module = fmt.Sprintf("%v%v", "../examples/modules/", instance.Module)
		} else {
			instance.Module = "dummy"
		}
		this.Relink()
		this.instanciate(instance)
		this.Relink()
	}
}
