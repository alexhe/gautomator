package flue

// rserver taken from the example of the libchan package

import (
	"crypto/tls"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"syscall"

	"github.com/docker/libchan"
	"github.com/docker/libchan/spdy"
)

// RemoteCommand is the received command parameters to execute locally and return
type RemoteCommand struct {
	Cmd        string
	Args       []string
	Stdin      io.Reader
	Stdout     io.WriteCloser
	Stderr     io.WriteCloser
	StatusChan libchan.Sender
}

// CommandResponse is the reponse struct to return to the client
type CommandResponse struct {
	Status int
}

func Rserver(proto *string, socket *string) {
	cert := os.Getenv("TLS_CERT")
	key := os.Getenv("TLS_KEY")

	var listener net.Listener
	if cert != "" && key != "" {
		tlsCert, err := tls.LoadX509KeyPair(cert, key)
		if err != nil {
			log.Fatal(err)
		}

		tlsConfig := &tls.Config{
			InsecureSkipVerify: true,
			Certificates:       []tls.Certificate{tlsCert},
		}

		listener, err = tls.Listen(*proto, *socket, tlsConfig)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		var err error
		listener, err = net.Listen(*proto, *socket)
		if err != nil {
			log.Fatal(err)
		}
	}

	tl, err := spdy.NewTransportListener(listener, spdy.NoAuthenticator)
	if err != nil {
		log.Fatal(err)
	}

	for {
		t, err := tl.AcceptTransport()
		if err != nil {
			log.Print(err)
			break
		}

		go func() {
			for {
				receiver, err := t.WaitReceiveChannel()
				if err != nil {
					log.Print(err)
					break
				}

				go func() {
					for {
						command := &RemoteCommand{}
						err := receiver.Receive(command)
						if err != nil {
							log.Print(err)
							break
						}

						cmd := exec.Command(command.Cmd, command.Args...)
						cmd.Stdout = command.Stdout
						cmd.Stderr = command.Stderr
						//cmd.Stdout = os.Stdout
						//cmd.Stderr = os.Stderr

						stdin, err := cmd.StdinPipe()
						if err != nil {
							log.Print(err)
							break
						}
						if err != nil {
							log.Fatal(err)
						}
						go func() {
							log.Println("Copying back to stdin")
							io.Copy(stdin, command.Stdin)
							//log.Println("Closing stdin")
							//stdin.Close()
						}()
						log.Println("Running the command")
						res := cmd.Run()
						log.Printf("Command finished with error: %v", err)
						log.Println("Done")
						//log.Println("Closing the Stdout")
						//command.Stdout.Close()
						//log.Println("Closing the Stderr")
						//command.Stderr.Close()
						log.Println("Assigning returnResult")
						returnResult := &CommandResponse{}
						if res != nil {
							log.Println("Res is not null")
							if exiterr, ok := res.(*exec.ExitError); ok {
								log.Println("Treating the error code")
								returnResult.Status = exiterr.Sys().(syscall.WaitStatus).ExitStatus()
								log.Println("We have the result")
							} else {
								log.Print(res)
								returnResult.Status = 10
							}
						}
						log.Println("Res is not null")
						err = command.StatusChan.Send(returnResult)
						log.Println("Finished")
						if err != nil {
							log.Print(err)
						}
					}
				}()
			}
		}()
	}
}
