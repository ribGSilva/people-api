package person

import (
	"context"
	"github.com/ribgsilva/person-api/persistence/v1/person"
)

func Search(ctx context.Context, request SearchRequest) ([]Person, error) {
	search, err := person.Search(ctx, person.SearchPerson{
		Tags: request.Tags,
	})
	if err != nil {
		return nil, err
	}

	persons := make([]Person, len(search))
	for i, p := range search {
		persons[i] = Person(p)
	}

	return persons, nil
}
