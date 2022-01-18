package main

import (
	"github.com/ASUIFT401ProjectGroup19/cam-backend-services/cmd/setup"
	"google.golang.org/grpc/reflection"
)

func main() {
	listener, gRPCServer, closeHandlers, httpGateway, err := setup.NewServer()
	if err != nil {
		panic(err)
	}
	defer closeHandlers()

	reflection.Register(gRPCServer)

	go func() {
		if err = gRPCServer.Serve(listener); err != nil {
			panic(err)
		}
	}()

	httpGateway()

}
