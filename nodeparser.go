package flue

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func ParseNode(filename *string) *TaskDefinition {
	taskDefJson, err := ioutil.ReadFile(*filename)

	if err != nil {
		fmt.Println("Err is ", err)
	}

	var taskDef *TaskDefinition
	json.Unmarshal(taskDefJson, taskDef)
	return taskDef
}
