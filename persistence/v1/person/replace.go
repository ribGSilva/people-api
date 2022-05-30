package person

import (
	"context"
	"fmt"
	"github.com/ribgsilva/person-api/sys"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

func Replace(ctx context.Context, p Person) error {
	collection := sys.S.Mongo.Database(sys.Configs.Mongo.Database).Collection("Person")

	hex, err := primitive.ObjectIDFromHex(p.Id)
	if err != nil {
		return fmt.Errorf("failure to parse id: %w", err)
	}

	tCtx, cancel := context.WithTimeout(ctx, sys.Configs.Mongo.OperationTimeout)
	defer cancel()

	_, err = collection.ReplaceOne(tCtx, bson.M{"_id": hex}, updatePerson{
		Id:              hex,
		Name:            p.Name,
		Type:            p.Type,
		Role:            p.Role,
		ContactDuration: p.ContactDuration,
		Tags:            p.Tags,
		UpdatedAt:       time.Now().UTC(),
	})
	if err != nil {
		return fmt.Errorf("failure to search id: %w", err)
	}

	return nil
}
