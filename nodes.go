package flue

// http://mholt.github.io/json-to-go/
type TaskDefinition []struct {
	Taskname string   `json:"taskName"`
	Module   string   `json:"module"`
	Args     []string `json:"args"`
	Hosts    []string `json:"hosts"`
}
