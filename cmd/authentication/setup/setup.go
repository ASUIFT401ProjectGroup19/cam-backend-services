package setup

import (
	"flag"
	"fmt"
	"net"

	"github.com/ASUIFT401ProjectGroup19/cam-backend-services/internal/apihandlers/authentication"
	"github.com/ASUIFT401ProjectGroup19/cam-backend-services/internal/camdb"
	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
)

const (
	envCfgKey = "service"
)

type Config struct {
	Auth *authentication.Config
	DB   *camdb.Config
	Port string `default:"10000"`
}

type Handler interface {
	Close()
	RegisterAPIServer(*grpc.Server)
}

func GetConfig() (*Config, error) {
	config := &Config{}

	flag.Usage = func() { // To print all accepted ENV vars when run with -h
		flag.PrintDefaults()
		_ = envconfig.Usage(envCfgKey, config)
	}
	flag.Parse()

	err := envconfig.Process(envCfgKey, config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func NewServer() (net.Listener, *grpc.Server, func(), error) {
	config, err := GetConfig()
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

	db, err := camdb.New(config.DB, logger)
	if err != nil {
		return nil, nil, nil, err
	}

	handlers := []Handler{
		authentication.New(config.Auth, db, logger),
	}

	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", config.Port))
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
