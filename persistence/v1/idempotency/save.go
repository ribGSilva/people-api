package idempotency

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ribgsilva/person-api/sys"
)

func Save(ctx context.Context, idempotency Idempotency) error {
	key, err := buildKey(search{
		Id:       idempotency.Id,
		Endpoint: idempotency.Endpoint,
	})
	if err != nil {
		return err
	}

	bytes, err := json.Marshal(idempotency)
	if err != nil {
		return fmt.Errorf("error parsing idempotency to be stored for key %s: %w", key, err)
	}

	tCtx, cancel := context.WithTimeout(ctx, sys.Configs.Redis.OperationTimeout)
	defer cancel()

	set, err := sys.S.Redis.SetNX(tCtx, key, string(bytes), sys.Configs.Redis.CacheTTL).Result()
	if err != nil {
		return fmt.Errorf("error searcing in cache for key %s: %w", key, err)
	}

	if !set {
		return fmt.Errorf("error saving in cache for key %s: %w", key, err)
	}

	return nil
}
