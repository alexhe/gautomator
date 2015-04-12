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
	log.Println("Server mode")
	proto := "tcp"
	socket := "localhost:4546"
	go gautomator.Rserver(&proto, &socket)

	if len(dotFiles) != 0 {
		log.Println("Client mode")

		taskStructure := gautomator.ParseDotFiles(dotFiles)
		// Parse the nodes.json and adapt the tasks
		taskInstances := gautomator.ParseNode(nodesFileJson)
		allSubTasks := make(map[int]*gautomator.TaskGraphStructure, 0)
		index := 0
		for _, taskInstance := range taskInstances {
			for _, node := range taskInstance.Hosts {
				subTasks := taskStructure.GetSubstructure(taskInstance.Taskname)
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
				}
			}
		}

		for _, subTask := range allSubTasks {
			//subTask.PrintAdjacencyMatrix()
			taskStructure = taskStructure.AugmentTaskStructure(subTask)
		}
		taskStructure.Relink()
		// Now, for each task, assign module, hosts and co...
		for _, task := range taskStructure.Tasks {
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
