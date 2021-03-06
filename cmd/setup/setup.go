package setup

import (
	"context"
	"flag"
	"fmt"
	commentHandler "github.com/ASUIFT401ProjectGroup19/cam-backend-services/internal/api/handlers/comment"
	feedHandler "github.com/ASUIFT401ProjectGroup19/cam-backend-services/internal/api/handlers/feed"
	galleryHandler "github.com/ASUIFT401ProjectGroup19/cam-backend-services/internal/api/handlers/gallery"
	identityHandler "github.com/ASUIFT401ProjectGroup19/cam-backend-services/internal/api/handlers/identity"
	postHandler "github.com/ASUIFT401ProjectGroup19/cam-backend-services/internal/api/handlers/post"
	subscriptionHandler "github.com/ASUIFT401ProjectGroup19/cam-backend-services/internal/api/handlers/subscription"
	"github.com/ASUIFT401ProjectGroup19/cam-backend-services/internal/api/middleware/interceptors/auth"
	"github.com/ASUIFT401ProjectGroup19/cam-backend-services/internal/api/middleware/interceptors/validation"
	storageAdapter "github.com/ASUIFT401ProjectGroup19/cam-backend-services/internal/core/adapters/persistence/cam"
	dbDriver "github.com/ASUIFT401ProjectGroup19/cam-backend-services/internal/core/adapters/persistence/cam/database/cam"
	sessionManager "github.com/ASUIFT401ProjectGroup19/cam-backend-services/internal/core/adapters/session"
	"github.com/ASUIFT401ProjectGroup19/cam-backend-services/internal/core/adapters/session/tokenmanager"
	commentServer "github.com/ASUIFT401ProjectGroup19/cam-backend-services/internal/core/servers/comment"
	feedServer "github.com/ASUIFT401ProjectGroup19/cam-backend-services/internal/core/servers/feed"
	galleryServer "github.com/ASUIFT401ProjectGroup19/cam-backend-services/internal/core/servers/gallery"
	identityServer "github.com/ASUIFT401ProjectGroup19/cam-backend-services/internal/core/servers/identity"
	postServer "github.com/ASUIFT401ProjectGroup19/cam-backend-services/internal/core/servers/post"
	subscriptionServer "github.com/ASUIFT401ProjectGroup19/cam-backend-services/internal/core/servers/subscription"
	commentV1 "github.com/ASUIFT401ProjectGroup19/cam-common/pkg/gen/proto/go/comment/v1"
	feedV1 "github.com/ASUIFT401ProjectGroup19/cam-common/pkg/gen/proto/go/feed/v1"
	galleryV1 "github.com/ASUIFT401ProjectGroup19/cam-common/pkg/gen/proto/go/gallery/v1"
	identityV1 "github.com/ASUIFT401ProjectGroup19/cam-common/pkg/gen/proto/go/identity/v1"
	postV1 "github.com/ASUIFT401ProjectGroup19/cam-common/pkg/gen/proto/go/post/v1"
	subscriptionV1 "github.com/ASUIFT401ProjectGroup19/cam-common/pkg/gen/proto/go/subscription/v1"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"
	"net/http"

	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
)

const (
	envCfgKey = "service"
)

type Config struct {
	DB           *dbDriver.Config
	Comment      *commentHandler.Config
	Feed         *feedHandler.Config
	Gallery      *galleryHandler.Config
	Identity     *identityHandler.Config
	Port         string `default:"10000"`
	Post         *postHandler.Config
	RestPort     string `default:"11000"`
	Subscription *subscriptionHandler.Config
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
		err := envconfig.Usage(envCfgKey, config)
		if err != nil {
			log.Fatal(err)
		}
	}
	flag.Parse()

	err := envconfig.Process(envCfgKey, config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func NewGRPCServer(config *Config) (net.Listener, *grpc.Server, func(), error) {
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

	databaseDriver, err := dbDriver.New(config.DB, logger)
	if err != nil {
		return nil, nil, nil, err
	}

	tm, err := tokenmanager.New(config.TokenManager)
	if err != nil {
		return nil, nil, nil, err
	}

	storage := storageAdapter.New(databaseDriver)

	session := sessionManager.New(storage, tm)

	handlers := []Handler{
		commentHandler.New(config.Comment, session, commentServer.New(storage), logger),
		feedHandler.New(config.Feed, session, feedServer.New(storage), logger),
		galleryHandler.New(config.Gallery, session, galleryServer.New(storage), logger),
		identityHandler.New(config.Identity, session, identityServer.New(storage), logger),
		postHandler.New(config.Post, session, postServer.New(storage), logger),
		subscriptionHandler.New(config.Subscription, session, subscriptionServer.New(storage), logger),
	}

	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", config.Port))
	if err != nil {
		return nil, nil, nil, err
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

	return listener, server, closeHandlers, nil
}

func NewHTTPServer(config *Config) (func() error, error) {
	mux := runtime.NewServeMux()

	cors := func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization, ResponseType")
			if r.Method == http.MethodOptions {
				return
			}
			h.ServeHTTP(w, r)
		})
	}

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	if err := commentV1.RegisterCommentServiceHandlerFromEndpoint(context.Background(), mux, fmt.Sprintf("localhost:%s", config.Port), opts); err != nil {
		return nil, err
	}
	if err := feedV1.RegisterFeedServiceHandlerFromEndpoint(context.Background(), mux, fmt.Sprintf("localhost:%s", config.Port), opts); err != nil {
		return nil, err
	}
	if err := galleryV1.RegisterGalleryServiceHandlerFromEndpoint(context.Background(), mux, fmt.Sprintf("localhost:%s", config.Port), opts); err != nil {
		return nil, err
	}
	if err := identityV1.RegisterIdentityServiceHandlerFromEndpoint(context.Background(), mux, fmt.Sprintf("localhost:%s", config.Port), opts); err != nil {
		return nil, err
	}
	if err := postV1.RegisterPostServiceHandlerFromEndpoint(context.Background(), mux, fmt.Sprintf("localhost:%s", config.Port), opts); err != nil {
		return nil, err
	}
	if err := subscriptionV1.RegisterSubscriptionServiceHandlerFromEndpoint(context.Background(), mux, fmt.Sprintf("localhost:%s", config.Port), opts); err != nil {
		return nil, err
	}
	gateway := func() error {
		if err := http.ListenAndServe(fmt.Sprintf(":%s", config.RestPort), cors(mux)); err != nil {
			return err
		}
		return nil
	}
	return gateway, nil
}
