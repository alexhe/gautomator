package flue

// This will convert the TaskGraphStructure into a format suitable for sigmajs
type SigmaNode struct {
	Id    int    // "id": "1",
	Labe  string //"label": "Node 1",
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
	Edges []SigmaEdge
	Nodes []SigmaNode
}

func GetSigmaTaskStructure(taskGraphStructure *TaskGraphStructure) *SigmaStructure {

}
