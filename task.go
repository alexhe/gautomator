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
	layer  string //the running layer
	module string
	args   []string
	status int // (2: schedule, 1:running, 0:done)
	rcode  int // The return code of the task (0 is ok)
}
