package schema

import (
	"context"
	"github.com/ribgsilva/person-api/sys"
)

func Create(ctx context.Context, _ []string) error {
	database := sys.S.Mongo.Database(sys.Configs.Mongo.Database)

	tCtx, cancel := context.WithTimeout(ctx, sys.Configs.Mongo.OperationTimeout)
	defer cancel()

	err := database.CreateCollection(tCtx, "Person")
	return err
}
