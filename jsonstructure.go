package gautomator

import (
	"fmt"
	"strconv"
)

// This will convert the TaskGraphStructure into a format suitable for sigmajs
type SigmaNode struct {
	Id    string  `json:"id"`    // "id": "1",
	Label string  `json:"label"` //"label": "Node 1",
	Color string  `json:"color"` //"color": "rgb(90,90,90)",
	Size  float64 `json:"size"`  //"size": 100,
	X     float64 `json:"x"`     //"x": 10,
	Y     float64 `json:"y"`     //"y": -10,
	Type  string  `json:"type"`  //"type": "tweetegy"

}

type SigmaEdge struct {
	Id     string `json:"id"`
	Source string `json:"source"`
	Target string `json:"target"`
	Type   string `json:"type"` //"type": "tweetegy"
}

type SigmaStructure struct {
	Nodes []*SigmaNode `json:"nodes"`
	Edges []*SigmaEdge `json:"edges"`
}

func (this *SigmaStructure) AddEdge(sigmaEdge *SigmaEdge) {
	this.Edges = append(this.Edges, sigmaEdge)
}
func (this *SigmaStructure) AddNode(sigmaNode *SigmaNode) {
	this.Nodes = append(this.Nodes, sigmaNode)
}

func NewSigmaStructure() *SigmaStructure {
	return &SigmaStructure{
		make([]*SigmaNode, 0),
		make([]*SigmaEdge, 0),
	}
}
func NewSigmaEdge() *SigmaEdge {
	return &SigmaEdge{
		string(-1),
		string(-1),
		string(-1),
		"curvedArrow",
	}
}

func NewSigmaNode() *SigmaNode {
	return &SigmaNode{
		"-1",
		"Default Node",
		"rgb(90,90,90)",
		100,
		0,
		0,
		"Def",
	}
}
func GetSigmaStructure(taskGraphStructure *TaskGraphStructure) *SigmaStructure {
	// First parse the taskGraphStructure.Tasks and create the node array
	var sigmaStructure *SigmaStructure
	sigmaStructure = NewSigmaStructure()
	for _, task := range taskGraphStructure.Tasks {
		sigmaNode := NewSigmaNode()
		sigmaNode.Id = strconv.Itoa(task.Id)
		sigmaNode.Label = fmt.Sprint(task.Name, ":", task.Node)
		sigmaStructure.AddNode(sigmaNode)
	}
	rowSize, colSize := taskGraphStructure.AdjacencyMatrix.Dims()
	edgeId := 1
	for r := 0; r < rowSize; r++ {
		for c := 0; c < colSize; c++ {
			// If there is a link, create the edge
			if taskGraphStructure.AdjacencyMatrix.At(r, c) != 0 {
				sigmaEdge := NewSigmaEdge()
				sigmaEdge.Id = strconv.Itoa(edgeId)
				sigmaEdge.Source = strconv.Itoa(r)
				sigmaEdge.Target = strconv.Itoa(c)
				edgeId += 1
				sigmaStructure.AddEdge(sigmaEdge)
			}
		}
	}
	return sigmaStructure
}
