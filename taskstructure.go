package flue

/*
A task is an atomic action performed by a module.

The **go** structure of a task in descripbed in the file **task.go**

a task is composed of:
- the layer of the task
- a module name
- an array of arguments
- a status of execution which can be
-- scheduled
-- running
-- done
- a return code
*/
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
}

// This is the structure corresponding to the "dot-graph" of a task list
// We store the nodes in a map
// The index is the source node
// The value is an array of strings containing the destination
type TaskGraphStructure struct {
        tasks []task
        deps  map[string][]string // A map index task1 will wait for task2, task3 and task4 to be completed
}
