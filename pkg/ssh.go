package pkg

import (
	"bufio"
	"fmt"
	"golang.org/x/crypto/ssh"
	command "grpc-demo/grpc"
	"io"
	"log"
	"net"
	sync "sync"
	"time"
)

const keys = `
xxx
`

type Stream struct {
	Stream  command.Command_CommandExecServer
}

// SSHClient ...
func (s Stream) SSHClient(hostname string, WG sync.WaitGroup) {
	signer, err := ssh.ParsePrivateKey([]byte(keys))
	if err != nil {
		log.Printf("unable to parse private key: %v", err)
	}

	config := ssh.ClientConfig{
		User: "root",
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
		Timeout: 5 * time.Second,
	}

	c, err := ssh.Dial("tcp", hostname+":22", &config)
	if err != nil {
		log.Println(err)
	}

	session, err := c.NewSession()
	if err != nil {
		log.Println(err)
	}
	defer session.Close()

	stderr, err := session.StderrPipe()
	if err != nil {
		log.Printf("opening stderr: %v", err)
	}
	stdout, err := session.StdoutPipe()
	if err != nil {
		log.Printf("opening stdout: %v", err)
	}

	if err := session.Start("ping -c 3 xxx.xxx.xxx.xxx"); err != nil {
		log.Println(err)
	}

	er := make(chan string)
	consumeStream := func(r io.Reader) {
		scan := bufio.NewScanner(r)
		scan.Split(bufio.ScanLines)
		for scan.Scan() {
			combo := scan.Text()
			er <- combo
		}
		if err := scan.Err(); err != nil {
			log.Println(err)
		}
	}


	go func() {
		if err := session.Wait(); err != nil {
			log.Printf("-------- %s --------", err)
		} else {
			time.Sleep(time.Duration(2)*time.Second)
			if err, ok := err.(*ssh.ExitError); ok {
				fmt.Println(err.ExitStatus())
			}
			if err, ok := err.(*ssh.ExitMissingError); !ok {
				fmt.Println(err.Error())
			} else {
				fmt.Println("------------ OK ---------")
			}
			er <- "The command is executed, please exit"
		}
	}()

	go consumeStream(stderr)
	go consumeStream(stdout)
	<- er

	defer WG.Done()

	for {
		combo := <- er
		if err := s.Stream.Send(&command.RunExecResponse{Hostname: hostname, Resp: combo}); err != nil {
			log.Println(err)
		}
	}
}

