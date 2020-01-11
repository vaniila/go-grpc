package main

import (
	"context"
	"fmt"

	"github.com/micro/cli"
	"github.com/micro/go-grpc"
	hello "github.com/micro/go-grpc/examples/greeter/server/proto/hello"
	"github.com/micro/go-micro"
)

var (
	// service to call
	serviceName string
)

func main() {
	service := grpc.NewService()

	service.Init(
		micro.Flags(cli.StringFlag{
			Name:        "service_name",
			Value:       "go.micro.srv.greeter",
			Destination: &serviceName,
		}),
	)

	cl := hello.NewSayClient(serviceName, service.Client())

	rsp, err := cl.Hello(context.TODO(), &hello.Request{
		Name: "John",
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(rsp.Msg)
}
