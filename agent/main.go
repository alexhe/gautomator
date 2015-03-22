package main

import (
	"github.com/owulveryck/flue"
	"log"
	"os"
)

func main() {
	//	flue.ParseTopology()
	//	flue.ParseNode()
	if len(os.Args) < 2 {
		log.Println("We are a server...")
		flue.Server("/tmp/mysocket.sock")
	} else {
		log.Println("We are a client...")
		command := &flue.RemoteCommandClient{
			Cmd:    os.Args[1],
			Args:   os.Args[2:],
			Stdin:  os.Stdin,
			Stdout: os.Stdout,
			Stderr: os.Stderr,
			//StatusChan: remoteSender,
		}
		flue.Client(command, "/tmp/mysocket.sock")
	}

}
