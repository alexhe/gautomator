package flue

// This will convert the TaskGraphStructure into a format suitable for sigmajs
type SigmaNode struct {
	Id    int    // "id": "1",
	Label string //"label": "Node 1",
	Color string //"color": "rgb(90,90,90)",
	Size  int    //"size": 100,
	X     int    //"x": 10,
	Y     int    //"y": -10,
	Type  string //"type": "tweetegy"

}

type SigmaEdge struct {
	Id     int
	Source int
	Target int
}

type SigmaStructure struct {
	Nodes []*SigmaNode
	Edges []*SigmaEdge
}

func NewSigmaStructure() *SigmaStructure {
	return &SigmaStructure{
		make([]*SigmaNode, 2),
		make([]*SigmaEdge, 1),
	}
}
func NewSigmaEdge() *SigmaEdge {
	return &SigmaEdge{
		-1,
		-1,
		-1,
	}
}

func NewSigmaNode() *SigmaNode {
	return &SigmaNode{
		-1,
		"Default Node",
		"rgb(90,90,90)",
		100,
		0,
		0,
		"Def",
	}
}

func GetSigmaTaskStructure(taskGraphStructure *TaskGraphStructure) *SigmaStructure {
	// First parse the taskGraphStructure.Tasks and create the node array
	return nil

}
