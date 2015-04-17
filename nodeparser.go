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
