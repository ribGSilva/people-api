package person

import (
	"context"
	"github.com/ribgsilva/person-api/persistence/v1/person"
)

func Create(ctx context.Context, createReq CreateRequest) (CreateResponse, error) {
	id, err := person.Save(ctx, person.NewPerson{
		Name:            createReq.Name,
		Type:            createReq.Type,
		Role:            createReq.Role,
		ContactDuration: createReq.ContactDuration,
		Tags:            createReq.Tags,
	})
	if err != nil {
		return CreateResponse{}, err
	}
	return CreateResponse{
		Id: id,
	}, nil
}
