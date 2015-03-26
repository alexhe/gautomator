package flue

import (
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"log"
	"os"
)

// This function will copy the current executable in remotePath
func MirrorMe(conn *ssh.Client, remotePath string) {
	// open an SFTP session over an existing ssh connection.
	sftp, err := sftp.NewClient(conn)
	if err != nil {
		log.Fatal(err)
	}
	defer sftp.Close()

	// walk a directory
	w := sftp.Walk(remotePath)
	for w.Step() {
		if w.Err() != nil {
			continue
		}
		log.Println(w.Path())
	}

	// The name of the executable
	filename := os.Args[1] // get command line first parameter

	f, err := sftp.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	// Read the content of the file
	dat, err := ioutil.ReadFile(filename)
	if _, err := f.Write(dat); err != nil {
		log.Fatal(err)
	}

	// check it's there
	fi, err := sftp.Lstat(filename)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(fi)
}
