package main

import (
	"github.com/ASUIFT401ProjectGroup19/cam-backend-services/cmd/setup"
	"google.golang.org/grpc/reflection"
)

func main() {
	config, err := setup.GetConfig()
	if err != nil {
		panic(err)
	}
	listener, gRPCServer, closeHandlers, err := setup.NewGRPCServer(config)
	defer closeHandlers()
	if err != nil {
		panic(err)
	}

	httpGateway, err := setup.NewHTTPServer(config)
	if err != nil {
		panic(err)
	}

	reflection.Register(gRPCServer)

	go func() {
		if err = gRPCServer.Serve(listener); err != nil {
			panic(err)
		}
	}()

	httpGateway()

}
