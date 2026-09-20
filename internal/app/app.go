// Package app configures and runs application.
package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/leabago/share-radio/adder/config"
	amqprpc "github.com/leabago/share-radio/adder/internal/controller/amqp_rpc"
	"github.com/leabago/share-radio/adder/internal/controller/grpc"
	grpcmw "github.com/leabago/share-radio/adder/internal/controller/grpc/middleware"
	"github.com/leabago/share-radio/adder/internal/controller/http_api"
	"github.com/leabago/share-radio/adder/internal/usecase/genre"
	"github.com/leabago/share-radio/adder/internal/usecase/language"
	"github.com/leabago/share-radio/adder/internal/usecase/station"

	persistGenreRepo "github.com/leabago/share-radio/adder/internal/repo/persistent/genre"
	persistLanguageRepo "github.com/leabago/share-radio/adder/internal/repo/persistent/language"
	persistStationRepo "github.com/leabago/share-radio/adder/internal/repo/persistent/station"

	"github.com/leabago/share-radio/adder/internal/usecase"
	"github.com/leabago/share-radio/adder/pkg/grpcserver"
	"github.com/leabago/share-radio/adder/pkg/httpserver"
	"github.com/leabago/share-radio/adder/pkg/jwt"
	"github.com/leabago/share-radio/adder/pkg/logger"
	"github.com/leabago/share-radio/adder/pkg/postgres"
	rmqRPCServer "github.com/leabago/share-radio/adder/pkg/rabbitmq/rmq_rpc/server"
	"github.com/leabago/share-radio/adder/pkg/tracing"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	pbgrpc "google.golang.org/grpc"
)

type useCases struct {
	station  usecase.Station
	genre    usecase.Genre
	language usecase.Language
}

type servers struct {
	rmq  *rmqRPCServer.Server
	grpc *grpcserver.Server
	http *httpserver.Server
}

func initUseCases(cfg *config.Config, pg *postgres.Postgres) useCases {

	stationRepo := persistStationRepo.New(pg)
	genreRepo := persistGenreRepo.New(pg)
	languageRepo := persistLanguageRepo.New(pg)

	return useCases{

		station:  station.New(cfg, stationRepo),
		genre:    genre.New(genreRepo),
		language: language.New(languageRepo),
	}
}

func initServers(cfg *config.Config, uc useCases, jwtManager *jwt.Manager, l logger.Interface) servers {
	// RabbitMQ RPC Server
	rmqRouter := amqprpc.NewRouter(jwtManager, l)

	rmqServer, err := rmqRPCServer.New(cfg.RMQ.URL, cfg.RMQ.ServerExchange, rmqRouter, l)
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - rmqServer - server.New: %w", err))
	}

	// gRPC Server
	grpcServer := grpcserver.New(
		l,
		grpcserver.Port(cfg.GRPC.Port),
		grpcserver.ServerOptions(
			pbgrpc.UnaryInterceptor(grpcmw.AuthInterceptor(jwtManager)),
			pbgrpc.StatsHandler(otelgrpc.NewServerHandler()),
		),
	)
	grpc.NewRouter(grpcServer.App, l)

	// HTTP Server
	httpServer := httpserver.New(l, httpserver.Port(cfg.HTTP.Port), httpserver.Prefork(cfg.HTTP.UsePreforkMode))
	http_api.NewRouter(httpServer.App, cfg, uc.station, uc.genre, uc.language, l)

	return servers{
		rmq:  rmqServer,
		grpc: grpcServer,
		http: httpServer,
	}
}

func (s *servers) startServers() {
	s.rmq.Start()
	s.grpc.Start()
	s.http.Start()

}

func (s *servers) waitForShutdown(l logger.Interface) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	var err error

	select {
	case sig := <-interrupt:
		l.Info("app - Run - signal: %s", sig.String())
	case err = <-s.http.Notify():
		l.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	case err = <-s.grpc.Notify():
		l.Error(fmt.Errorf("app - Run - grpcServer.Notify: %w", err))
	case err = <-s.rmq.Notify():
		l.Error(fmt.Errorf("app - Run - rmqServer.Notify: %w", err))
	}

	s.shutdownServers(l)
}

func (s *servers) shutdownServers(l logger.Interface) {
	if err := s.http.Shutdown(); err != nil {
		l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}

	if err := s.grpc.Shutdown(); err != nil {
		l.Error(fmt.Errorf("app - Run - grpcServer.Shutdown: %w", err))
	}

	if err := s.rmq.Shutdown(); err != nil {
		l.Error(fmt.Errorf("app - Run - rmqServer.Shutdown: %w", err))
	}

}

// Run creates objects via constructors.
func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)

	ctx := context.Background()

	// Tracing
	shutdownTracing, err := tracing.New(ctx, tracing.Config{
		Enabled:     cfg.Tracing.Enabled,
		ServiceName: cfg.App.Name,
		Version:     cfg.App.Version,
		Endpoint:    cfg.Tracing.OTLPEndpoint,
		Insecure:    cfg.Tracing.OTLPInsecure,
		SampleRate:  cfg.Tracing.SampleRate,
	})
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - tracing.New: %w", err))
	}
	defer func() {
		if err := shutdownTracing(ctx); err != nil {
			l.Error(fmt.Errorf("app - Run - shutdownTracing: %w", err))
		}
	}()

	// Repository
	pg, err := postgres.New(cfg.PG.URL, postgres.MaxPoolSize(cfg.PG.PoolMax))
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err))
	}
	defer pg.Close()

	// JWT
	jwtManager := jwt.New(cfg.JWT.Secret, cfg.JWT.TokenExpiry)

	uc := initUseCases(cfg, pg)
	s := initServers(cfg, uc, jwtManager, l)
	s.startServers()
	s.waitForShutdown(l)
}
