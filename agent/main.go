package main

import (
	"flag"
	"github.com/owulveryck/gautomator"
	"log"
	"net/http"
	"sync"
)

func cleanup() {
	log.Println("cleanup")
}

func main() {

	//Parsing the dot
	var dotFiles []string
	var nodesFileJson = flag.String("nodes", "", "json file for node definition")
	flag.Parse()
	dotFiles = flag.Args()
	if len(dotFiles) == 0 {
		log.Println("Server mode")
		proto := "tcp"
		socket := "localhost:4546"
		gautomator.Rserver(&proto, &socket)
	} else {
		log.Println("Client mode")

		taskStructure := gautomator.ParseDotFiles(dotFiles)
		// Parse the nodes.json and adapt the tasks
		taskStructure.PrintDot()
		nodeStructure := gautomator.ParseNode(nodesFileJson)
		//var allSubTasks []*gautomator.TaskGraphStructure
		allSubTasks := make(map[int]*gautomator.TaskGraphStructure, 0)
		index := 0
		for _, nodeDef := range nodeStructure {
			for _, node := range nodeDef.Hosts {
				subTasks := taskStructure.GetSubstructure(nodeDef.Taskname)
				// If there is subtask
				if subTasks != nil {
					for i, _ := range subTasks.Tasks {
						subTasks.Tasks[i].Node = node
					}
					allSubTasks[index] = subTasks
					index += 1
				} else {
					// TODO Duplicate a single task
				}
			}
		}
		/*
		for _, subTask := range allSubTasks {
			//subTask.PrintAdjacencyMatrix()
			taskStructure = taskStructure.AugmentTaskStructure(subTask)
		}
		taskStructure.Relink()
		*/
		//taskStructure.PrintAdjacencyMatrix()
		// Entering the workers area
		var wg sync.WaitGroup
		doneChan := make(chan *gautomator.Task)

		// For each task, launch a goroutine
		for taskIndex, _ := range taskStructure.Tasks {
			go gautomator.Runner(taskStructure.Tasks[taskIndex], doneChan, &wg)
			wg.Add(1)
		}
		go gautomator.Advertize(taskStructure, doneChan)

		// This is the web displa
		router := gautomator.NewRouter(taskStructure)

		go log.Fatal(http.ListenAndServe(":8080", router))

		// Wait for all the runner(s) to be finished
		wg.Wait()
	}
}
