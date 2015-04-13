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

func (this *TaskGraphStructure) InstanciateTaskStructure(taskInstances TaskDefinition) {
		allSubTasks := make(map[int]*TaskGraphStructure, 0)
		index := 0
		for _, taskInstance := range taskInstances {
			for _, node := range taskInstance.Hosts {
				subTasks := this.GetSubstructure(taskInstance.Taskname)
				// If there is subtask
				if subTasks != nil {
					for i, _ := range subTasks.Tasks {
						log.Printf("Setting node %v on task %v (%v)", node, subTasks.Tasks[i].Name, i)
						subTasks.Tasks[i].Node = node
					}
					allSubTasks[index] = subTasks
					index += 1
				} else {
					// TODO Duplicate a single task
					// Get the id of the task to duplicate
					log.Printf("DEBUG: Duplicating %v",taskInstance.Taskname)
					//newId := this.DuplicateTask(taskInstance.Taskname)
					//this.Tasks[newId].Node = node
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
					log.Printf("DEBUG module %v (%v)", taskInstances[task.Name].Module, task.Name)
					task.Module = taskInstances[task.Name].Module
				}
				if taskInstances[task.Name].Args != nil {
					log.Printf("DEBUG Args %v (%v)", taskInstances[task.Name].Args, task.Name)
					task.Args = taskInstances[task.Name].Args
				}
			}
		}



}
