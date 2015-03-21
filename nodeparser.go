package flue

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func ParseNode() {
	rolesJson, err := ioutil.ReadFile("books/nodes.json")

	if err != nil {
		fmt.Println("Err is ", err)
	}

	var roles Roles
	json.Unmarshal(rolesJson, &roles)
	fmt.Println(roles)
}
