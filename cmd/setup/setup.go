package setup

import (
	"context"
	"flag"
	"fmt"
	"github.com/ASUIFT401ProjectGroup19/cam-backend-services/internal/servers/identity"
	"github.com/ASUIFT401ProjectGroup19/cam-backend-services/internal/servers/post"
	identityGatewayv1 "github.com/ASUIFT401ProjectGroup19/cam-common/pkg/gen/proto/go/identity/v1"
	postGatewayv1 "github.com/ASUIFT401ProjectGroup19/cam-common/pkg/gen/proto/go/post/v1"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/credentials/insecure"
	"net"
	"net/http"

	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"

	storageAdapter "github.com/ASUIFT401ProjectGroup19/cam-backend-services/internal/adapters/persistence/cam"
	authHandler "github.com/ASUIFT401ProjectGroup19/cam-backend-services/internal/apihandlers/identity"
	postHandler "github.com/ASUIFT401ProjectGroup19/cam-backend-services/internal/apihandlers/post"
	camDriver "github.com/ASUIFT401ProjectGroup19/cam-backend-services/internal/database/cam"
	"github.com/ASUIFT401ProjectGroup19/cam-backend-services/internal/middleware/interceptors/auth"
	"github.com/ASUIFT401ProjectGroup19/cam-backend-services/internal/middleware/interceptors/validation"
	"github.com/ASUIFT401ProjectGroup19/cam-backend-services/internal/middleware/tokenmanager"
)

const (
	envCfgKey = "service"
)

type Config struct {
	Auth         *authHandler.Config
	DB           *camDriver.Config
	Port         string `default:"10000"`
	Post         *postHandler.Config
	RestPort     string `default:"11000"`
	TokenManager *tokenmanager.Config
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

func NewServer() (net.Listener, *grpc.Server, func(), func(), error) {
	config, err := GetConfig()
	if err != nil {
		return nil, nil, nil, nil, err
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
		return nil, nil, nil, nil, err
	}

	databaseDriver, err := camDriver.New(config.DB, logger)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	session, err := tokenmanager.New(config.TokenManager)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	storage := storageAdapter.New(databaseDriver)

	identityServer := identity.New(session, storage)
	postServer := post.New(session, storage)

	handlers := []Handler{
		authHandler.New(config.Auth, identityServer, logger),
		postHandler.New(config.Post, postServer, logger),
	}

	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", config.Port))
	if err != nil {
		return nil, nil, nil, nil, err
	}

	authInterceptor := auth.New(session)
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

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	if err := identityGatewayv1.RegisterIdentityServiceHandlerFromEndpoint(context.Background(), mux, fmt.Sprintf("localhost:%s", config.Port), opts); err != nil {
		panic(err)
	}
	if err := postGatewayv1.RegisterPostServiceHandlerFromEndpoint(context.Background(), mux, fmt.Sprintf("localhost:%s", config.Port), opts); err != nil {
		panic(err)
	}
	gateway := func() {
		if err := http.ListenAndServe(fmt.Sprintf(":%s", config.RestPort), mux); err != nil {
			panic(err)
		}
	}

	return listener, server, closeHandlers, gateway, nil
}