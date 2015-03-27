package flue

import (
    "log"
	"github.com/awalterschulze/gographviz"
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
    taskDstExists := false // the flag to check if a task already exist
    taskSrcExists := false // the flag to check if a task already exist
    for _, aTask := range this.Tasks {
	if aTask != nil { 
	    //log.Println("Current task is",aTask.Name)
	    // If the task exists, add src as a dependency
	    if aTask.Name == dst {
		taskDstExists = true
		aTask.Deps = AppendString(aTask.Deps,src)
	    }
	    if aTask.Name == src {
		taskSrcExists = true
	    }
	} 
    }
    // If the task does not exists, create it and add it to the structure
    if taskSrcExists == false {
	aTask := NewTask()
	aTask.Name = src
	this.Tasks = AppendTask(this.Tasks,aTask)
    }
    if taskDstExists == false {
	aTask := NewTask()
	aTask.Name = dst
	aTask.Deps = AppendString(aTask.Deps,src)
	this.Tasks = AppendTask(this.Tasks,aTask)
    }
}
func (this *TaskGraphStructure) AddEdge(src, dst string, directed bool, attrs map[string]string) {
	this.AddPortEdge(src, "", dst, "", directed, attrs)
}
func (this *TaskGraphStructure) AddNode(parentGraph string, name string, attrs map[string]string) {
    log.Printf("parentGraph: %v, name: %v",parentGraph,name)
    for key, value := range attrs {
	log.Printf("Arg: %v, Value:%v",key,value)
    }

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
