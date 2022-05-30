package person

import (
	"context"
	"fmt"
	"github.com/ribgsilva/person-api/sys"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func Search(ctx context.Context, search SearchPerson) ([]Person, error) {
	collection := sys.S.Mongo.Database(sys.Configs.Mongo.Database).Collection("Person")

	tCtx, cancel := context.WithTimeout(ctx, sys.Configs.Mongo.OperationTimeout)
	defer cancel()

	cur, err := collection.Find(tCtx, bson.M{"tags": bson.M{"$in": search.Tags}})
	if err != nil {
		if err != mongo.ErrNoDocuments {
			return nil, fmt.Errorf("failure to search documents: %w", err)
		} else {
			return nil, nil
		}
	}
	defer cur.Close(ctx)

	persons := make([]Person, 0)
	for cur.Next(ctx) {
		var p person
		if err := cur.Decode(&p); err != nil {
			return nil, fmt.Errorf("failure to parse person search: %w", err)
		}
		persons = append(persons, Person{
			Id:              p.Id.Hex(),
			Name:            p.Name,
			Type:            p.Type,
			Role:            p.Role,
			ContactDuration: p.ContactDuration,
			Tags:            p.Tags,
			UpdatedAt:       p.UpdatedAt,
			CreatedAt:       p.CreatedAt,
		})
	}

	return persons, nil

}
