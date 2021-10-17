package main

import (
	"github.com/ASUIFT401ProjectGroup19/cam-backend-services/cmd/authentication/setup"
)

func main() {
	listener, server, closeHandlers, err := setup.NewServer()
	if err != nil {
		panic(err)
	}

	if err = server.Serve(listener); err != nil {
		panic(err)
	}

	closeHandlers()
}
