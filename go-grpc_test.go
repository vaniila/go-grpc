package grpc

import (
	"context"
	"crypto/tls"
	"sync"
	"testing"

	hello "github.com/micro/go-grpc/examples/greeter/server/proto/hello"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry/mock"
	mls "github.com/micro/misc/lib/tls"
)

type testHandler struct{}

func (t *testHandler) Hello(ctx context.Context, req *hello.Request, rsp *hello.Response) error {
	rsp.Msg = "Hello " + req.Name
	return nil
}

func TestGRPCService(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// create mock registry
	r := mock.NewRegistry()

	// create GRPC service
	service := NewService(
		micro.Name("test.service"),
		micro.Registry(r),
		micro.AfterStart(func() error {
			wg.Done()
			return nil
		}),
		micro.Context(ctx),
	)

	// register test handler
	hello.RegisterSayHandler(service.Server(), &testHandler{})

	// run service
	go func() {
		if err := service.Run(); err != nil {
			t.Fatal(err)
		}
	}()

	// wait for start
	wg.Wait()

	// create client
	say := hello.NewSayClient("test.service", service.Client())

	// call service
	rsp, err := say.Hello(context.Background(), &hello.Request{
		Name: "John",
	})
	if err != nil {
		t.Fatal(err)
	}

	// check message
	if rsp.Msg != "Hello John" {
		t.Fatalf("unexpected response %s", rsp.Msg)
	}
}

func TestGRPCFunction(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// create service
	fn := NewFunction(
		micro.Name("test.function"),
		micro.Registry(mock.NewRegistry()),
		micro.AfterStart(func() error {
			wg.Done()
			return nil
		}),
		micro.Context(ctx),
	)

	// register test handler
	hello.RegisterSayHandler(fn.Server(), &testHandler{})

	// run service
	go fn.Run()

	// wait for start
	wg.Wait()

	// create client
	say := hello.NewSayClient("test.function", fn.Client())

	// call service
	rsp, err := say.Hello(context.Background(), &hello.Request{
		Name: "John",
	})
	if err != nil {
		t.Fatal(err)
	}

	// check message
	if rsp.Msg != "Hello John" {
		t.Fatalf("unexpected response %s", rsp.Msg)
	}
}

func TestGRPCTLSService(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// create mock registry
	r := mock.NewRegistry()

	// create cert
	cert, err := mls.Certificate("test.service")
	if err != nil {
		t.Fatal(err)
	}
	config := &tls.Config{
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: true,
	}

	// create GRPC service
	service := NewService(
		micro.Name("test.service"),
		micro.Registry(r),
		micro.AfterStart(func() error {
			wg.Done()
			return nil
		}),
		micro.Context(ctx),
		// set TLS config
		WithTLS(config),
	)

	// register test handler
	hello.RegisterSayHandler(service.Server(), &testHandler{})

	// run service
	go func() {
		if err := service.Run(); err != nil {
			t.Fatal(err)
		}
	}()

	// wait for start
	wg.Wait()

	// create client
	say := hello.NewSayClient("test.service", service.Client())

	// call service
	rsp, err := say.Hello(context.Background(), &hello.Request{
		Name: "John",
	})
	if err != nil {
		t.Fatal(err)
	}

	// check message
	if rsp.Msg != "Hello John" {
		t.Fatalf("unexpected response %s", rsp.Msg)
	}
}
