package person

import (
	"context"
	"github.com/ribgsilva/person-api/persistence/v1/person"
)

func Get(ctx context.Context, id string) (Person, error) {
	find, err := person.Find(ctx, id)
	if err != nil {
		return Person{}, err
	}
	if find.Id == "" {
		return Person{}, nil
	}
	return Person(find), nil
}
