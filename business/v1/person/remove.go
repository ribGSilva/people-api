package person

import (
	"context"
	"github.com/ribgsilva/person-api/persistence/v1/person"
)

func Remove(ctx context.Context, id string) (Person, error) {
	del, err := person.Delete(ctx, id)
	if err != nil {
		return Person{}, err
	}
	if del.Id == "" {
		return Person{}, nil
	}
	return Person(del), nil
}
