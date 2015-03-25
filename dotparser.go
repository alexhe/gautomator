package flue

import (
	"github.com/awalterschulze/gographviz"
)

// We store the nodes in a map
// The index is the source node
// The value is an array of strings containing the destination
type TopologyGraphStructure struct {
	allTheTasks []string
	waiter      map[string][]string // A map index task1 will wait for task2, task3 and task4 to be completed
}

func NewTopologyGraphStructure() *TopologyGraphStructure {
	return &TopologyGraphStructure{
		make([]string, 0),
		make(map[string][]string),
	}
}

func appendTask(slice []string, task string) []string {
	for _, ele := range slice {
		if ele == task {
			return slice
		}
	}
	return append(slice, task)
}

// Compsite literal ?
// http://golang.org/ref/spec#Composite_literals
//func Newmap[string][]string() *map[string][]string {
//	return &map[string][]string{
//	return make(map[string][]string)
// Here wi shall make the array
//	}
//}

func (this *TopologyGraphStructure) SetStrict(strict bool) {}
func (this *TopologyGraphStructure) SetDir(directed bool)  {}
func (this *TopologyGraphStructure) SetName(name string)   {}
func (this *TopologyGraphStructure) AddPortEdge(src, srcPort, dst, dstPort string, directed bool, attrs map[string]string) {
	this.allTheTasks = appendTask(this.allTheTasks, src)
	this.waiter[dst] = append(this.waiter[dst], src)
}
func (this *TopologyGraphStructure) AddEdge(src, dst string, directed bool, attrs map[string]string) {
	this.AddPortEdge(src, "", dst, "", directed, attrs)
}
func (this *TopologyGraphStructure) AddNode(parentGraph string, name string, attrs map[string]string) {
}
func (this *TopologyGraphStructure) AddAttr(parentGraph string, field, value string) {}
func (this *TopologyGraphStructure) AddSubGraph(parentGraph string, name string, attrs map[string]string) {
}
func (this *TopologyGraphStructure) String() string { return "" }

func ParseTopology(topologyDot []byte) *TopologyGraphStructure {

	parsed, err := gographviz.Parse(topologyDot)
	if err != nil {
		panic(err)
	}
	// Display the graph
	//fmt.Println(parsed)
	var topology *TopologyGraphStructure
	topology = NewTopologyGraphStructure()
	gographviz.Analyse(parsed, topology)
	//fmt.Println(topology.role["Ref2"][0])
	//fmt.Println(topology.role["Ref1"][1])
	return topology
}
