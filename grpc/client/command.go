package client

import (
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	command "grpc-demo/grpc"
	"io"
	"log"
	"time"
)

type response struct {
	Hostname string `json:"hostname"`
	Resp string	`json:"resp"`
}

func NewSSHResponse(h, res string) *response {
	r := new(response)
	r.Hostname = h
	r.Resp = res
	return r
}

const (
	address = "localhost:50052"
)

func HandleExec(hosts []string, cmd string) ([]*response, error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Panicln(err)
	}
	defer conn.Close()

	c := command.NewCommandClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*100)
	defer cancel()

	stream, err := c.CommandExec(ctx)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if err := stream.Send(&command.RunExec{Hostname: hosts, Cmd: cmd}); err != nil {
		log.Println(err)
	}
	if err := stream.CloseSend(); err != nil {
		log.Println(err)
	}
	
	var res []*response

	for {
		es, err := stream.Recv()
		if err == io.EOF {
			log.Println(err)
			break
		}
		if err != nil {
			log.Println(err)
		}

		fmt.Println(es.Hostname, es.Resp)
		r := NewSSHResponse(es.Hostname, es.Resp)
		res = append(res, r)
	}
	return res, err
}