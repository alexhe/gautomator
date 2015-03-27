package flue

import (
	"github.com/awalterschulze/gographviz"
)

func AppendTask(slice []string, task string) []string {
	for _, ele := range slice {
		if ele == task {
			return slice
		}
	}
	return append(slice, task)
}

func (this *TaskGraphStructure) SetStrict(strict bool) {}
func (this *TaskGraphStructure) SetDir(directed bool)  {}
func (this *TaskGraphStructure) SetName(name string)   {}
func (this *TaskGraphStructure) AddPortEdge(src, srcPort, dst, dstPort string, directed bool, attrs map[string]string) {
	this.AllTheTasks = AppendTask(this.AllTheTasks, src)
	this.waiter[dst] = append(this.waiter[dst], src)
}
func (this *TaskGraphStructure) AddEdge(src, dst string, directed bool, attrs map[string]string) {
	this.AddPortEdge(src, "", dst, "", directed, attrs)
}
func (this *TaskGraphStructure) AddNode(parentGraph string, name string, attrs map[string]string) {
}
func (this *TaskGraphStructure) AddAttr(parentGraph string, field, value string) {}
func (this *TaskGraphStructure) AddSubGraph(parentGraph string, name string, attrs map[string]string) {
}
func (this *TaskGraphStructure) String() string { return "" }

func ParseTask(topologyDot []byte) *TaskGraphStructure {

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
