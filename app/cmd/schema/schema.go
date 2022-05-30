package schema

import (
	"context"
	"fmt"
	"github.com/ribgsilva/person-api/persistence/v1/schema"
	"github.com/ribgsilva/person-api/platform/env"
	"github.com/ribgsilva/person-api/sys"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.uber.org/zap"
)

func ListCommands() {
	println("Schema Commands")
	println("\tcreate\t\t\t- Creates the schema")
	println("\tdelete\t\t\t- Deletes the schema")
	println("\thelp\t\t\t- Print the commands available")
}

func Run(options []string) {
	if len(options) == 0 {
		ListCommands()
		return
	}
	// empty logger
	log := zap.NewNop().Sugar()
	if err := initVars(log); err != nil {
		println("error:", err)
	}
	defer func() {
		sdCtx, sdCancel := context.WithTimeout(context.Background(), sys.Configs.Mongo.DisconnectTimeout)
		defer sdCancel()
		if err := sys.S.Mongo.Disconnect(sdCtx); err != nil {
			log.Error(err)
		}
	}()
	switch options[0] {
	case "create":
		println("creating schema")
		if err := schema.Create(context.Background(), options[1:]); err != nil {
			println("failed to create schema:", err.Error())
		} else {
			println("created schema")
		}
	case "delete":
		println("deleting schema")
		if err := schema.Delete(context.Background(), options[1:]); err != nil {
			println("failed to delete schema:", err.Error())
		} else {
			println("deleted schema")
		}
	case "help":
		fallthrough
	default:
		ListCommands()
	}
}

func initVars(log *zap.SugaredLogger) error {
	sys.Configs.Mongo.Database = env.OrDefault(log, "MONGO_DATABASE", "Person")
	sys.Configs.Mongo.ConnectionURL = env.OrDefault(log, "MONGO_CONNECTION_URL", "mongodb://person_app:person_app_pass@localhost:27017")
	sys.Configs.Mongo.ConnectionTimeout = env.DurationDefault(log, "MONGO_CONNECTION_TIMEOUT", "10s")
	sys.Configs.Mongo.DisconnectTimeout = env.DurationDefault(log, "MONGO_DISCONNECT_TIMEOUT", "5s")
	sys.Configs.Mongo.OperationTimeout = env.DurationDefault(log, "MONGO_OPERATION_TIMEOUT", "10s")
	sys.Configs.Mongo.PingTimeout = env.DurationDefault(log, "MONGO_PING_TIMEOUT", "2s")

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
}
