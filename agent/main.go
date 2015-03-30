package main

import (
	//	"fmt"
	//	"github.com/nu7hatch/gouuid"
	"flag"
	"github.com/owulveryck/flue"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	//	"os"
	//	    "syscall"
	//	"os/signal"
)

func cleanup() {
	log.Println("cleanup")
}

func main() {

	/*
		// Catching the interrupt
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		signal.Notify(c, syscall.SIGTERM)
		go func() {
		    <-c
		    cleanup()
		    os.Exit(1)
		}()
	*/
	//Parsing the dot
	dotFile := flag.String("dot", "", "The dot file")

	flag.Parse()
	if *dotFile == "" {
		log.Println("Server mode")
		proto := "tcp"
		socket := "localhost:4546"
		flue.Rserver(&proto, &socket)
	} else {
		log.Println("Client mode")

		var topologyDot []byte
		topologyDot, err := ioutil.ReadFile(*dotFile)
		if err != nil {
			log.Panic("Cannot read file")
		}

		log.Println("Parsing...")
		taskStructure := flue.ParseTasks(topologyDot)
		// How many tasks

		var wg sync.WaitGroup
		taskStructureChan := make(chan *flue.TaskGraphStructure)
		doneChan := make(chan *flue.Task)
		for i, task := range taskStructure.Tasks {
			if task != nil {
				go flue.Runner(taskStructure, task, taskStructureChan, doneChan, &wg)
				wg.Add(1)
				log.Printf("=> Tache %v: %s", i, task.Name)
				for j, dep := range task.Deps {
					log.Printf("==> Deps[%v]: %v", j, dep)
				}
			}
		}
		go flue.Advertize(taskStructure, taskStructureChan, doneChan)
		router := flue.NewRouter(*taskStructure)

		go log.Fatal(http.ListenAndServe(":8080", router))
		wg.Wait()
	}
}
