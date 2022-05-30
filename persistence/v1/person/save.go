package person

import (
	"context"
	"github.com/ribgsilva/person-api/sys"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

func Save(ctx context.Context, np NewPerson) (string, error) {
	collection := sys.S.Mongo.Database(sys.Configs.Mongo.Database).Collection("Person")

	now := time.Now().UTC()
	newP := newPerson{
		Name:            np.Name,
		Type:            np.Type,
		Role:            np.Role,
		ContactDuration: np.ContactDuration,
		Tags:            np.Tags,
		UpdatedAt:       now,
		CreatedAt:       now,
	}

	tCtx, cancel := context.WithTimeout(ctx, sys.Configs.Mongo.OperationTimeout)
	defer cancel()

	res, err := collection.InsertOne(tCtx, newP)
	if err != nil {
		return "", err
	}
	return res.InsertedID.(primitive.ObjectID).Hex(), nil
}
