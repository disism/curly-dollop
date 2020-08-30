package main

import (
	"fmt"
	"google.golang.org/grpc"
	command "grpc-demo/grpc"
	"grpc-demo/pkg"
	"io"
	"log"
	"net"
	"sync"
)

type server struct {
	command.CommandServer
}

// CommandExec ...
func (s *server) CommandExec(stream command.Command_CommandExecServer) error {
	var WG sync.WaitGroup
	for {
		ex, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		ns := new(pkg.Stream)
		ns.Stream = stream

		WG.Add(len(ex.Hostname))
		for _, item := range ex.Hostname {
			go ns.SSHClient(item, WG)
		}
		WG.Wait()
	}
}

func main() {
	listen, err := net.Listen("tcp", ":50052")
	if err != nil{
		log.Println(err)
	}
	s := grpc.NewServer()

	command.RegisterCommandServer(s, &server{})
	fmt.Println("gRpc Server: 50052")
	if err := s.Serve(listen); err != nil {
		log.Println(err)
	}
}