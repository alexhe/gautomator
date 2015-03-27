package flue

// A task is an action executed by a module
type Task struct {
        Name   string //the task name
        Node   string // The node name
        Module string
        Args   []string
        Status int // 0: not run yet
        // 1: running
        // 2: finished with success
        // 3: finished with error
        ReturnCode int // The return code of the task (0 is ok)
        Deps  []string // A map index task1 will wait for task2, task3 and task4 to be completed
}

// This is the structure corresponding to the "dot-graph" of a task list
// We store the nodes in a map
// The index is the source node
// The value is an array of strings containing the destination
type TaskGraphStructure struct {
        Tasks []*Task
}

func NewTask() *Task {
    return &Task{
	"null",
	"localhost",
	"dummy",
	make([]string, 0),
	0,
	0,
	make([]string, 0),
    }

}
func NewTaskGraphStructure() *TaskGraphStructure {
	return &TaskGraphStructure{
		make([]*Task,0),
	}
}


