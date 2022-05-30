package person

import (
	"context"
	"github.com/ribgsilva/person-api/persistence/v1/person"
)

func Replace(ctx context.Context, id string, request UpdateRequest) (Person, error) {
	find, err := person.Find(ctx, id)
	if err != nil {
		return Person{}, err
	}
	if find.Id == "" {
		return Person{}, nil
	}
	up := person.Person{
		Id:              find.Id,
		Name:            request.Name,
		Type:            request.Type,
		Role:            request.Role,
		ContactDuration: request.ContactDuration,
		Tags:            request.Tags,
		CreatedAt:       find.CreatedAt,
	}
	if err := person.Replace(ctx, up); err != nil {
		return Person{}, err
	}
	return Person(up), nil
}
