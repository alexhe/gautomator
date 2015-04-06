package flue

import (
	"github.com/awalterschulze/gographviz"
	"github.com/gonum/matrix/mat64"
	"log"
)

func AppendTask(slice []*Task, task *Task) []*Task {
	for _, element := range slice {
		if element == task {
			return slice
		}
	}
	return append(slice, task)
}
func AppendString(slice []string, task string) []string {
	for _, element := range slice {
		if element == task {
			return slice
		}
	}
	return append(slice, task)
}

func (this *TaskGraphStructure) SetStrict(strict bool) {}
func (this *TaskGraphStructure) SetDir(directed bool)  {}
func (this *TaskGraphStructure) SetName(name string)   {}
func (this *TaskGraphStructure) AddPortEdge(src, srcPort, dst, dstPort string, directed bool, attrs map[string]string) {
	lastIndex := len(this.Tasks)
	// Find the index of the task src and the index of the task dst
	srcTaskId := -1
	dstTaskId := -1
	increment := 0 // The number of new lines and cols needed
	for taskId, taskObject := range this.Tasks {
		if taskObject != nil {
			//log.Println("Current task is",taskObject.Name)
			// If the task exists, add src as a dependency
			if taskObject.Name == dst {
				dstTaskId = taskId
			}
			if taskObject.Name == src {
				srcTaskId = taskId
			}
		}
	}
	// If the task does not exists, create it and add it to the structure
	if srcTaskId == -1 {
		taskObject := NewTask()
		taskObject.Name = src
		this.Tasks[lastIndex] = taskObject
		increment += 1
		srcTaskId = lastIndex
		taskObject.Id = srcTaskId
		lastIndex += 1
	}
	// If the task does not exists, create it and add it to the structure
	if dstTaskId == -1 {
		taskObject := NewTask()
		taskObject.Name = dst
		this.Tasks[lastIndex] = taskObject
		increment += 1
		dstTaskId = lastIndex
		taskObject.Id = dstTaskId
		lastIndex += 1
	}
	// If the size of the increment is not null
	// use Grow...
	if increment > 0 {
		this.DegreeMatrix = mat64.DenseCopyOf(this.DegreeMatrix.Grow(increment, increment))
		this.AdjacencyMatrix = mat64.DenseCopyOf(this.AdjacencyMatrix.Grow(increment, increment))
	}
	// Now fill the matrix
	this.DegreeMatrix.Set(dstTaskId, dstTaskId, this.DegreeMatrix.At(dstTaskId, dstTaskId)+1)
	this.DegreeMatrix.Set(srcTaskId, srcTaskId, this.DegreeMatrix.At(srcTaskId, srcTaskId)+1)
	this.AdjacencyMatrix.Set(srcTaskId, dstTaskId, 1)
}

//TODO: write a merge function to merge two structures:
// -> Concat the task map
// -> merge the matrix
func (this *TaskGraphStructure) AddEdge(src, dst string, directed bool, attrs map[string]string) {
	this.AddPortEdge(src, "", dst, "", directed, attrs)
}
func (this *TaskGraphStructure) AddNode(parentGraph string, name string, attrs map[string]string) {
	for _, taskObject := range this.Tasks {
		if taskObject != nil && taskObject.Name == name {
			return
		}
	}
	id := len(this.Tasks)
	taskObject := NewTask()
	taskObject.Name = name
	taskObject.Id = id
	taskObject.Origin = parentGraph
	this.Tasks[id] = taskObject
	this.DegreeMatrix = mat64.DenseCopyOf(this.DegreeMatrix.Grow(1, 1))
	this.AdjacencyMatrix = mat64.DenseCopyOf(this.AdjacencyMatrix.Grow(1, 1))
}
func (this *TaskGraphStructure) AddAttr(parentGraph string, field, value string) {}
func (this *TaskGraphStructure) AddSubGraph(parentGraph string, name string, attrs map[string]string) {
}
func (this *TaskGraphStructure) String() string { return "" }

func ParseTasks(topologyDot []byte) *TaskGraphStructure {

	parsed, err := gographviz.Parse(topologyDot)
	if err != nil {
		panic(err)
	}
	// Display the graph
	//fmt.Println(parsed)
	var topology *TaskGraphStructure
	topology = NewTaskGraphStructure()
	gographviz.Analyse(parsed, topology)
	//fmt.Println(topology.role["Ref2"][0])
	//fmt.Println(topology.role["Ref1"][1])
	return topology
}
