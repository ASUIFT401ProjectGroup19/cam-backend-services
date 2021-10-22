package setup

import (
	"fmt"
	"net"

	"go.uber.org/zap/zapcore"

	authenticationAPI "github.com/ASUIFT401ProjectGroup19/cam-backend-services/internal/apihandlers/authentication"

	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

const (
	envCfgKey = "service"
)

type Config struct {
	authenticationAPI.Config
	Port string `default:"10000"`
}

type Handler interface {
	Close()
	RegisterAPIServer(*grpc.Server)
}

func NewServer() (net.Listener, *grpc.Server, func(), error) {
	conf := &Config{}
	err := envconfig.Process(envCfgKey, conf)
	if err != nil {
		return nil, nil, nil, err
	}

	logger, err := zap.Config{
		Level:    zap.NewAtomicLevelAt(zap.DebugLevel),
		Encoding: "json",
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey: "message",
		},
		OutputPaths: []string{"stdout"},
	}.Build()
	if err != nil {
		return nil, nil, nil, err
	}

	handlers := []Handler{
		authenticationAPI.New(logger),
	}

	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", conf.Port))
	if err != nil {
		return nil, nil, nil, err
	}

	server := grpc.NewServer()

	for _, handler := range handlers {
		handler.RegisterAPIServer(server)
	}

	closeHandlers := func() {
		for _, handler := range handlers {
			handler.Close()
		}
	}

	return listener, server, closeHandlers, nil
}
