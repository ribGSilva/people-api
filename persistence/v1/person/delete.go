package person

import (
	"context"
	"fmt"
	"github.com/ribgsilva/person-api/sys"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func Delete(ctx context.Context, id string) (Person, error) {
	collection := sys.S.Mongo.Database(sys.Configs.Mongo.Database).Collection("Person")

	tCtx, cancel := context.WithTimeout(ctx, sys.Configs.Mongo.OperationTimeout)
	defer cancel()

	hex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return Person{}, fmt.Errorf("failure to parse id: %w", err)
	}
	one := collection.FindOneAndDelete(tCtx, bson.M{"_id": hex})

	if err := one.Err(); err != nil {
		if err != mongo.ErrNoDocuments {
			return Person{}, fmt.Errorf("failure to delete id: %w", err)
		} else {
			return Person{}, nil
		}
	}

	var p person
	if err := one.Decode(&p); err != nil {
		return Person{}, fmt.Errorf("failure to parse delete person by id %s: %w", id, err)
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
