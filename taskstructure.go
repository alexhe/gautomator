package flue

// A task is an action executed by a module
type task struct {
        name   string //the task name
        node   string // The node name
        module string
        args   []string
        status int // 0: not run yet
        // 1: running
        // 2: finished with success
        // 3: finished with error
        returnCode int // The return code of the task (0 is ok)
        deps  []string // A map index task1 will wait for task2, task3 and task4 to be completed
}

// This is the structure corresponding to the "dot-graph" of a task list
// We store the nodes in a map
// The index is the source node
// The value is an array of strings containing the destination
type TaskGraphStructure struct {
        tasks []task
}

func NewTaskGraphStructure() *TaskGraphStructure {
	return &TaskGraphStructure{
		make(map[string]int),
		make(map[string][]string),
	}
}


