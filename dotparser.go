package main

import (
	"fmt"
	"github.com/awalterschulze/gographviz"
	"io/ioutil"
)

// We store the nodes in a map
// The index is the source node
// The value is an array of strings containing the destination
type TopologyGraphStructure struct {
	role map[string][]string
}

func NewTopologyGraphStructure() *TopologyGraphStructure {
	return &TopologyGraphStructure{
		make(map[string][]string),
	}
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
	// If SourceName already exists, add destination to the structure
	// else, add a new entry in the structure
	this.role[src] = append(this.role[src], dst)
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

func ParseTopology() *TopologyGraphStructure {
	// Testing the DOT parsing...
	topologyDot, err := ioutil.ReadFile("books/topology.dot")
	if err != nil {
		fmt.Println("Err is ", err)
	}

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
