package flue

import (
	"crypto/tls"
	"github.com/docker/libchan"
	"github.com/docker/libchan/spdy"
	"io"
	"log"
	"net"
	"os"
)

// Command is the run parameters to be executed remotely
type Command struct {
	Cmd        string
	Args       []string
	Stdin      io.Writer
	Stdout     io.Reader
	Stderr     io.Reader
	StatusChan libchan.Sender
}

func Client(task *Task, proto *string, socket *string) int {
	var client net.Conn
	var err error
	if os.Getenv("USE_TLS") != "" {
		client, err = tls.Dial(*proto, *socket, &tls.Config{InsecureSkipVerify: true})
	} else {
		client, err = net.Dial(*proto, *socket)
	}
	if err != nil {
		log.Fatal(err)
	}
	transport, err := spdy.NewClientTransport(client)
	if err != nil {
		log.Fatal(err)
	}
	sender, err := transport.NewSendChannel()
	if err != nil {
		log.Fatal(err)
	}
	//receiver, remoteSender := libchan.Pipe()
	_, remoteSender := libchan.Pipe()
	command := &Command{
		Cmd:        task.Module,
		Args:       task.Args[:],
		Stdin:      os.Stdin,
		Stdout:     os.Stdout,
		Stderr:     os.Stderr,
		StatusChan: remoteSender,
	}
	log.Println("Sending command")
	err = sender.Send(command)
	if err != nil {
		log.Fatal(err)
	}
	//sender.Close()
	response := &CommandResponse{}
	log.Println("Receiving response")
	//err = receiver.Receive(response)
	log.Println("Received")
	if err != nil {
		log.Fatal(err)
	}
	command.StatusChan.Close()
	sender.Close()
	remoteSender.Close()
	client.Close()
	return response.Status
	//	os.Exit(response.Status)
}
