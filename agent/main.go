package main

import (
	"github.com/nu7hatch/gouuid"
	"github.com/owulveryck/flue"
	"log"
	"os"
)

func main() {
	// Testing the DOT parsing...
	//	topologyDot, err := ioutil.ReadFile("test.dot")
	topologyDot := []byte(`
digraph layer3Tasks {
	start -> purge;
	purge -> installProduit1;
	purge -> installProduit2;
	installProduit1 -> startAll;
	installProduit2 -> startAll;
	startAll -> end;	

}
`)
	flue.ParseTopology(topologyDot)
	//	flue.ParseNode()

	if len(os.Args) < 2 {
		uuid, _ := uuid.NewV4()
		log.Println("We are a server, uuid is: ", string(uuid[:]))
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
