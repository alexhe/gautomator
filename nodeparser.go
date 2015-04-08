package gautomator

import (
	"encoding/json"
	"fmt"
	"log"
	"io/ioutil"
)

func ParseNode(filename *string) TaskDefs {
	taskDefJson, err := ioutil.ReadFile(*filename)

	if err != nil {
		fmt.Println("Err is ", err)
	}

	var taskDef TaskDefs
	err = json.Unmarshal(taskDefJson, &taskDef)
	if err != nil {
	    log.Panic(err)
	}
	return taskDef
}
