package gautomator

// http://mholt.github.io/json-to-go/
type TaskInstance struct {
	Taskname string   `json:"taskName"`
	Module   string   `json:"module"`
	Args     []string `json:"args"`
	Hosts    []string `json:"hosts"`
}

type taskDefs []TaskInstance

type TaskDefinition map[string]TaskInstance
