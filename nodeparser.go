package gautomator

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

func ParseNode(filename *string) TaskDefinition {
	taskDefJson, err := ioutil.ReadFile(*filename)

	if err != nil {
		fmt.Println("Err is ", err)
	}

	var taskDef taskDefs
	taskDefinition := make(map[string]TaskInstance, 0)

	err = json.Unmarshal(taskDefJson, &taskDef)
	if err != nil {
		log.Panic(err)
	}
	for _, task := range taskDef {
		taskDefinition[task.Taskname] = task
	}
	return taskDefinition
}

// TODO: this function is ugly and buggy...
// Need to rethink it and recode it in a more mathematical and elegant way
func (this *TaskGraphStructure) InstanciateTaskStructure(taskInstances TaskDefinition) {
	allSubTasks := make(map[int]*TaskGraphStructure, 0)
	index := 0
	for _, taskInstance := range taskInstances {
		for _, node := range taskInstance.Hosts {
			doNotDuplicate := false
			for _, task := range this.Tasks {
				if task.Origin == taskInstance.Taskname && task.Node == "null" {
					// Setting the node to node
					task.Node = node
					doNotDuplicate = true
				}
			}

			if doNotDuplicate != true {
				subTasks := this.GetSubstructure(taskInstance.Taskname)
				// If there is subtask
				if subTasks != nil {
					// TODO,  if a subtask exists with dummy, set it the hostname and no not add it
					for i, _ := range subTasks.Tasks {
						subTasks.Tasks[i].Node = node
					}
					allSubTasks[index] = subTasks
					index += 1
				} else {
					for _, task := range this.Tasks {
						if task.Name == taskInstance.Taskname {
							if task.Node == "null" || task.Node == node {
								task.Node = node
							} else {
								newIds := this.DuplicateTask(taskInstance.Taskname)
								for _, newId := range newIds {
									if newId != -1 {
										this.Tasks[newId].Node = node
									}
								}
							}
						}
					}
				}
			}
		}
	}

	for _, subTask := range allSubTasks {
		//subTask.PrintAdjacencyMatrix()
		this = this.AugmentTaskStructure(subTask)
	}
	this.Relink()
	// Now, for each task, assign module, hosts and co...
	for _, task := range this.Tasks {
		if _, ok := taskInstances[task.Name]; ok {
			if taskInstances[task.Name].Module != "" {
				task.Module = taskInstances[task.Name].Module
			}
			if taskInstances[task.Name].Args != nil {
				task.Args = taskInstances[task.Name].Args
			}
		}
	}

}
