package setup

import (
	"flag"
	"fmt"
	"net"

	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"

	"github.com/ASUIFT401ProjectGroup19/cam-backend-services/internal/adapters/persistence/camadapter"
	authHandler "github.com/ASUIFT401ProjectGroup19/cam-backend-services/internal/apihandlers/authentication"
	postHandler "github.com/ASUIFT401ProjectGroup19/cam-backend-services/internal/apihandlers/post"
	"github.com/ASUIFT401ProjectGroup19/cam-backend-services/internal/database/cam"
	"github.com/ASUIFT401ProjectGroup19/cam-backend-services/internal/middleware/interceptors/auth"
	"github.com/ASUIFT401ProjectGroup19/cam-backend-services/internal/middleware/interceptors/validation"
	"github.com/ASUIFT401ProjectGroup19/cam-backend-services/internal/middleware/tokenmanager"
)

const (
	envCfgKey = "service"
)

type Config struct {
	Auth         *authHandler.Config
	DB           *cam.Config
	TokenManager *tokenmanager.Config
	Port         string `default:"10000"`
	Post         *postHandler.Config
}

type Handler interface {
	Close()
	GetProtectedRPCs() []string
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

	databaseDriver, err := cam.New(config.DB, logger)
	if err != nil {
		return nil, nil, nil, err
	}

	tm, err := tokenmanager.New(config.TokenManager)
	if err != nil {
		return nil, nil, nil, err
	}

	persistenceAdapter := camadapter.New(databaseDriver)

	handlers := []Handler{
		authHandler.New(config.Auth, persistenceAdapter, logger, tm),
		postHandler.New(config.Post, persistenceAdapter, logger),
	}

	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", config.Port))
	if err != nil {
		return nil, nil, nil, err
	}

	authInterceptor := auth.New(tm)
	validationInterceptor := validation.New()

	unaryInterceptors := []grpc.UnaryServerInterceptor{
		authInterceptor.Unary(),
		validationInterceptor.Unary(),
	}

	streamInterceptors := []grpc.StreamServerInterceptor{
		authInterceptor.Stream(),
		validationInterceptor.Stream(),
	}

	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(unaryInterceptors...),
		grpc.ChainStreamInterceptor(streamInterceptors...),
	)

	for _, handler := range handlers {
		handler.RegisterAPIServer(server)
		authInterceptor.RegisterProtectedRoutes(handler.GetProtectedRPCs())
	}

	closeHandlers := func() {
		for _, handler := range handlers {
			handler.Close()
		}
	}

	return listener, server, closeHandlers, nil
}
