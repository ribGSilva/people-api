package person

import (
	"context"
	"fmt"
	"github.com/ribgsilva/person-api/sys"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func Find(ctx context.Context, id string) (Person, error) {
	collection := sys.S.Mongo.Database(sys.Configs.Mongo.Database).Collection("Person")

	tCtx, cancel := context.WithTimeout(ctx, sys.Configs.Mongo.OperationTimeout)
	defer cancel()

	hex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return Person{}, fmt.Errorf("failure to parse id: %w", err)
	}

	var p person
	err = collection.FindOne(tCtx, bson.M{"_id": hex}).Decode(&p)

	if err == mongo.ErrNoDocuments {
		return Person{}, nil
	} else if err != nil {
		return Person{}, fmt.Errorf("failure to search id: %w", err)
	}

	return Person{
		Id:              p.Id.Hex(),
		Name:            p.Name,
		Type:            p.Type,
		Role:            p.Role,
		ContactDuration: p.ContactDuration,
		Tags:            p.Tags,
		UpdatedAt:       p.UpdatedAt,
		CreatedAt:       p.CreatedAt,
	}, nil
}
