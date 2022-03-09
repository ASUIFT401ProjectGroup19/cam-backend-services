package main

import (
	"github.com/ASUIFT401ProjectGroup19/cam-backend-services/cmd/setup"
	"google.golang.org/grpc/reflection"
	"log"
)

func main() {
	config, err := setup.GetConfig()
	if err != nil {
		log.Fatal(err)
	}

	listener, gRPCServer, closeHandlers, err := setup.NewGRPCServer(config)
	defer closeHandlers()
	if err != nil {
		log.Fatal(err)
	}

	httpGateway, err := setup.NewHTTPServer(config)
	if err != nil {
		log.Fatal(err)
	}

	reflection.Register(gRPCServer)

	go func() {
		if err = gRPCServer.Serve(listener); err != nil {
			log.Fatal(err)
		}
	}()

	if err = httpGateway(); err != nil {
		log.Fatal(err)
	}

}
