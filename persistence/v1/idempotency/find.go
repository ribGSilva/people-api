package idempotency

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ribgsilva/person-api/sys"
)

func Find(ctx context.Context, id, endpoint string) (Idempotency, error) {
	key, err := buildKey(search{
		Id:       id,
		Endpoint: endpoint,
	})
	if err != nil {
		return Idempotency{}, err
	}

	tCtx, cancel := context.WithTimeout(ctx, sys.Configs.Redis.OperationTimeout)
	defer cancel()

	get, err := sys.S.Redis.Get(tCtx, key).Result()
	if err != nil {
		return Idempotency{}, fmt.Errorf("error searcing in cache for key %s: %w", key, err)
	}

	var idempotency Idempotency
	if err := json.Unmarshal([]byte(get), &idempotency); err != nil {
		return Idempotency{}, fmt.Errorf("error parsing cached response for key %s: %w", key, err)
	}

	return idempotency, nil
}
