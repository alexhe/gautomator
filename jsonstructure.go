package gautomator

import (
	"fmt"
	"strconv"
)

// This will convert the TaskGraphStructure into a format suitable for sigmajs
type jsonNode struct {
	Id    string  `json:"id"`    // "id": "1",
	Label string  `json:"label"` //"label": "Node 1",
	Color string  `json:"color"` //"color": "rgb(90,90,90)",
	Size  float64 `json:"size"`  //"size": 100,
	X     float64 `json:"x"`     //"x": 10,
	Y     float64 `json:"y"`     //"y": -10,
	Type  string  `json:"type"`  //"type": "tweetegy"

}

type jsonEdge struct {
	Id     string `json:"id"`
	Source string `json:"source"`
	Target string `json:"target"`
	Type   string `json:"type"` //"type": "tweetegy"
}

type jsonStructure struct {
	Nodes []*jsonNode `json:"nodes"`
	Edges []*jsonEdge `json:"edges"`
}

func (this *jsonStructure) AddEdge(sigmaEdge *jsonEdge) {
	this.Edges = append(this.Edges, sigmaEdge)
}
func (this *jsonStructure) AddNode(sigmaNode *jsonNode) {
	this.Nodes = append(this.Nodes, sigmaNode)
}

func NewjsonStructure() *jsonStructure {
	return &jsonStructure{
		make([]*jsonNode, 0),
		make([]*jsonEdge, 0),
	}
}
func NewjsonEdge() *jsonEdge {
	return &jsonEdge{
		string(-1),
		string(-1),
		string(-1),
		"curvedArrow",
	}
}

func NewjsonNode() *jsonNode {
	return &jsonNode{
		"-1",
		"Default Node",
		"rgb(90,90,90)",
		100,
		0,
		0,
		"Def",
	}
}
func GetjsonStructure(taskGraphStructure *TaskGraphStructure) *jsonStructure {
	// First parse the taskGraphStructure.Tasks and create the node array
	var sigmaStructure *jsonStructure
	//sigmaStructure = NewjsonStructure()
	for _, task := range taskGraphStructure.Tasks {
		sigmaNode := NewjsonNode()
		sigmaNode.Id = strconv.Itoa(task.Id)
		sigmaNode.Label = fmt.Sprint(task.Name, ":", task.Node)
		sigmaStructure.AddNode(sigmaNode)
	}
	return sigmaStructure
}
