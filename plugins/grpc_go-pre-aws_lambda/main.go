package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	coprocess "github.com/asoorm/tyk-mw-grpcgo-lambda/proto/go"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

const (
	listenNetwork = "unix"          // can be unix || tcp
	listenAddress = "/tmp/foo.sock" // can be path to unix socket || ip/port combo e.g. 127.0.0.1:3333
	awsRegion     = endpoints.EuWest2RegionID
)

func main() {

	go func() {
		<-handleSIGINTKILL()

		log.Println("received termination signal")

		if err := os.Remove(listenAddress); err != nil {
			log.Println(errors.Wrap(err, "unable to unbind, delete sock file manually"))
			os.Exit(1)
			return
		}

		os.Exit(0)
	}()

	listener, err := net.Listen(listenNetwork, listenAddress)
	if err != nil {
		log.Println(errors.Wrap(err, "error opening listener"))
		os.Exit(1)
		return
	}

	log.Printf("gRPC server listening on %s\n", listenAddress)

	server := grpc.NewServer()
	coprocess.RegisterDispatcherServer(server, NewDispatcher())

	log.Println(errors.Wrap(server.Serve(listener), "unable to serve"))
}

func handleSIGINTKILL() chan os.Signal {
	sig := make(chan os.Signal, 1)

	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	return sig
}
