package main

import (
	"github.com/TykTechnologies/tyk/coprocess"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"net"
)

// Dispatcher implementation
type Dispatcher struct{}

// Dispatch will be called on every request:
func (d *Dispatcher) Dispatch(ctx context.Context, object *coprocess.Object) (*coprocess.Object, error) {
	log.Println("Receiving object:", object)

	// We dispatch the object based on the hook name (as specified in the manifest file), these functions are in hooks.go:
	switch object.HookName {
	case "MyPreHook":
		log.Println("MyPreHook is called!")
		return MyPreHook(object)
	}

	log.Println("Unknown hook: ", object.HookName)

	return object, nil
}

// DispatchEvent will be called when a Tyk event is triggered:
func (d *Dispatcher) DispatchEvent(ctx context.Context, event *coprocess.Event) (*coprocess.EventReply, error) {
	return &coprocess.EventReply{}, nil
}

const (
	ListenAddress = ":5555"
)

// MyPreHook performs a header injection:
func MyPreHook(object *coprocess.Object) (*coprocess.Object, error) {
	object.Request.SetHeaders = map[string]string{
		"Myheader": "Myvalue",
	}

	return object, nil
}

func main() {
	lis, err := net.Listen("tcp", ListenAddress)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Printf("starting grpc server on %v", ListenAddress)
	s := grpc.NewServer()
	coprocess.RegisterDispatcherServer(s, &Dispatcher{})
	s.Serve(lis)
}
