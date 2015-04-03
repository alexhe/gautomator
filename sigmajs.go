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
func GetSigmaStructure(taskGraphStructure *TaskGraphStructure) *SigmaStructure {
	// First parse the taskGraphStructure.Tasks and create the node array
	var sigmaStructure *SigmaStructure
	sigmaStructure = NewSigmaStructure()
	for _, task := range taskGraphStructure.Tasks {
		sigmaNode := NewSigmaNode()
		sigmaNode.Id = task.Id
		sigmaNode.Label = task.Name
		sigmaStructure.AddNode(sigmaNode)
	}
	rowSize, colSize := taskGraphStructure.AdjacencyMatrix.Dims()
	edgeId := 1
	for r := 0; r < rowSize; r++ {
		for c := 0; c < colSize; c++ {
			// If there is a link, create the edge
			if taskGraphStructure.AdjacencyMatrix.At(r, c) != 0 {
				sigmaEdge := NewSigmaEdge()
				sigmaEdge.Id = edgeId
				sigmaEdge.Source = r
				sigmaEdge.Target = c
				edgeId += 1
				sigmaStructure.AddEdge(sigmaEdge)
			}
		}
	}
	return sigmaStructure
}
