package main

// http://mholt.github.io/json-to-go/
type Roles struct {
	Roles []struct {
		Type  string `json:"Type"`
		Nodes []struct {
			Hostname   string `json:"hostname"`
			User       string `json:"user"`
			SSH        string `json:"ssh"`
			Datasource []struct {
				Driver string `json:"driver"`
				Name   string `json:"name"`
			} `json:"datasource"`
		} `json:"nodes"`
	} `json:"Roles"`
}
