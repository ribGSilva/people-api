package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/newrelic/go-agent/v3/integrations/nrgin"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/ribgsilva/person-api/app/api/docs"
	"github.com/ribgsilva/person-api/app/api/handlers"
	"github.com/ribgsilva/person-api/platform/env"
	"github.com/ribgsilva/person-api/platform/logger"
	"github.com/ribgsilva/person-api/sys"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.uber.org/automaxprocs/maxprocs"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

// @title Person API
// @version 1.0
// @description Service to store handle person data.
// @contact.name Gabriel Ribeiro Silva
func main() {
	log, err := logger.New("Person-API")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer func(log *zap.SugaredLogger) {
		_ = log.Sync()
	}(log)

	if err := run(log); err != nil {
		log.Errorw("startup", "ERROR", err)
		_ = log.Sync()
		os.Exit(1)
	}
}

func run(log *zap.SugaredLogger) error {

	// =======================================================================================================
	// Setup max procs
	if _, err := maxprocs.Set(); err != nil {
		return fmt.Errorf("maxprocs: %w", err)
	}
	log.Infow("startup", "GOMAXPROCS", runtime.GOMAXPROCS(0))

	// =======================================================================================================
	// Setup configs
	sys.Configs.Http.Port = env.OrDefault(log, "HTTP_PORT", "8080")
	sys.Configs.Http.ReadTimeout = env.DurationDefault(log, "HTTP_SHUTDOWN_TIMEOUT", "5s")
	sys.Configs.Http.IdleTimeout = env.DurationDefault(log, "HTTP_SHUTDOWN_TIMEOUT", "120s")
	sys.Configs.Http.WriteTimeout = env.DurationDefault(log, "HTTP_SHUTDOWN_TIMEOUT", "10s")
	sys.Configs.Http.ShutdownTimeout = env.DurationDefault(log, "HTTP_SHUTDOWN_TIMEOUT", "60s")
	sys.Configs.Swagger.Protocol = env.OrDefault(log, "SWAGGER_PROTOCOL", "http")
	sys.Configs.Swagger.Host = env.OrDefault(log, "SWAGGER_HOST", "localhost:"+sys.Configs.Http.Port)
	sys.Configs.Mongo.Database = env.OrDefault(log, "MONGO_DATABASE", "Person")
	sys.Configs.Mongo.ConnectionURL = env.OrDefault(log, "MONGO_CONNECTION_URL", "mongodb://person_app:person_app_pass@localhost:27017")
	sys.Configs.Mongo.ConnectionTimeout = env.DurationDefault(log, "MONGO_CONNECTION_TIMEOUT", "20s")
	sys.Configs.Mongo.DisconnectTimeout = env.DurationDefault(log, "MONGO_DISCONNECT_TIMEOUT", "5s")
	sys.Configs.Mongo.PingTimeout = env.DurationDefault(log, "MONGO_PING_TIMEOUT", "5s")
	sys.Configs.Mongo.OperationTimeout = env.DurationDefault(log, "MONGO_OPERATION_TIMEOUT", "10s")
	sys.Configs.Redis.ConnectionURL = env.OrDefault(log, "REDIS_CONNECTION_URL", "localhost:6379")
	sys.Configs.Redis.User = env.OrDefault(log, "REDIS_USER", "")
	sys.Configs.Redis.Pass = env.OrDefault(log, "REDIS_PASS", "")
	sys.Configs.Redis.PingTimeout = env.DurationDefault(log, "REDIS_PING_TIMEOUT", "2s")
	sys.Configs.Redis.OperationTimeout = env.DurationDefault(log, "REDIS_PING_TIMEOUT", "10s")
	sys.Configs.Redis.CacheTTL = env.DurationDefault(log, "REDIS_CACHE_TTL", "24h")
	sys.Configs.NewRelic.AppName = env.OrDefault(log, "NEW_RELIC_APP_NAME", "person-api")
	sys.Configs.NewRelic.Licence = env.OrDefault(log, "NEW_RELIC_LICENCE", "")
	sys.Configs.NewRelic.Enabled = env.BoolDefault(log, "NEW_RELIC_ENABLED", "f")
	sys.Configs.NewRelic.ConnectionTimeout = env.DurationDefault(log, "NEW_RELIC_CONNECTION_TIMEOUT", "10s")
	sys.Configs.NewRelic.ShutdownTimeout = env.DurationDefault(log, "NEW_RELIC_SHUTDOWN_TIMEOUT", "10s")
	sys.Configs.IdempotencyEnabled = env.BoolDefault(log, "IDEMPOTENCY_ENABLED", "f")

	// =======================================================================================================
	// Setup static resources

	// logger
	sys.S.Log = log

	// mongo
	// doing in a func, so I can use defer to cancel the contexts
	if err := func() error {
		mongoCtx, mongoCancel := context.WithTimeout(context.Background(), sys.Configs.Mongo.ConnectionTimeout)
		defer mongoCancel()
		if client, err := mongo.Connect(mongoCtx, options.Client().ApplyURI(sys.Configs.Mongo.ConnectionURL)); err != nil {
			return fmt.Errorf("could not connect to mongo: %w", err)
		} else {
			pingCtx, pingCancel := context.WithTimeout(context.Background(), sys.Configs.Mongo.PingTimeout)
			defer pingCancel()
			if err = client.Ping(pingCtx, readpref.Primary()); err != nil {
				return fmt.Errorf("could not connect to mongo: %w", err)
			}
			sys.S.Mongo = client
		}
		return nil
	}(); err != nil {
		return err
	}
	defer func() {
		sdCtx, sdCancel := context.WithTimeout(context.Background(), sys.Configs.Mongo.DisconnectTimeout)
		defer sdCancel()
		if err := sys.S.Mongo.Disconnect(sdCtx); err != nil {
			log.Error(err)
		}
	}()

	// redis
	// doing in a func, so I can use defer to cancel the contexts
	if err := func() error {
		rdb := redis.NewClient(&redis.Options{
			Addr:     sys.Configs.Redis.ConnectionURL,
			Username: sys.Configs.Redis.User,
			Password: sys.Configs.Redis.Pass,
		})
		rdsCtx, rdsCancel := context.WithTimeout(context.Background(), sys.Configs.Redis.PingTimeout)
		defer rdsCancel()
		if err := rdb.Ping(rdsCtx).Err(); err != nil {
			return fmt.Errorf("could not connect to redis: %w", err)
		}
		sys.S.Redis = rdb
		return nil
	}(); err != nil {
		return err
	}
	defer func() {
		_ = sys.S.Redis.Close()
	}()

	// =======================================================================================================
	// NR

	nrApp, err := newrelic.NewApplication(
		newrelic.ConfigAppName(sys.Configs.NewRelic.AppName),
		newrelic.ConfigLicense(sys.Configs.NewRelic.Licence),
		newrelic.ConfigEnabled(sys.Configs.NewRelic.Enabled),
	)
	if err != nil {
		return err
	}
	if err := nrApp.WaitForConnection(sys.Configs.NewRelic.ConnectionTimeout); err != nil {
		return err
	}
	defer nrApp.Shutdown(sys.Configs.NewRelic.ShutdownTimeout)

	// =======================================================================================================
	// Router configuration

	router := gin.New()
	router.Use(gin.LoggerWithConfig(gin.LoggerConfig{
		SkipPaths: []string{"/v1/healthcheck"},
	}), gin.Recovery(), nrgin.Middleware(nrApp))

	handlers.MapHandlers(router)

	docs.SwaggerInfo.Host = sys.Configs.Swagger.Host
	url := ginSwagger.URL(fmt.Sprintf("%s://%s/swagger/doc.json", sys.Configs.Swagger.Protocol, sys.Configs.Swagger.Host))
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	// =======================================================================================================
	// App start and shutdown

	svr := &http.Server{
		Addr:         fmt.Sprintf(":%s", sys.Configs.Http.Port),
		Handler:      router,
		ReadTimeout:  sys.Configs.Http.ReadTimeout,
		WriteTimeout: sys.Configs.Http.WriteTimeout,
		IdleTimeout:  sys.Configs.Http.IdleTimeout,
	}

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	serverErrors := make(chan error, 1)
	go func() {
		log.Info("started http server")
		serverErrors <- svr.ListenAndServe()
	}()

	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)
	case sig := <-shutdown:
		log.Infow("shutdown", "status", "shutdown started", "signal", sig)
		defer log.Infow("shutdown", "status", "shutdown complete", "signal", sig)

		ctx, cancel := context.WithTimeout(context.Background(), sys.Configs.Http.ShutdownTimeout)
		defer cancel()

		if err := svr.Shutdown(ctx); err != nil {
			_ = svr.Close()
			return fmt.Errorf("could not stop server gracefully: %w", err)
		}
	}
	return nil
}
