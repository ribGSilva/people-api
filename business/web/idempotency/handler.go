package idempotency

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/ribgsilva/person-api/persistence/v1/idempotency"
	"github.com/ribgsilva/person-api/platform/web/handler"
	"github.com/ribgsilva/person-api/sys"
	"net/http"
)

type Handler struct {
	Endpoint string
	F        handler.HandleFunc
}

func (h Handler) Handle(ctx *gin.Context) handler.Result {

	id := ctx.GetHeader("X-Idempotency-Key")

	if !sys.Configs.IdempotencyEnabled || id == "" {
		return h.F(ctx)
	}

	err := idempotency.Save(ctx, idempotency.Idempotency{
		Id:       id,
		Endpoint: h.Endpoint,
	})

	if err != nil {
		cached, err := idempotency.Find(ctx, id, h.Endpoint)
		if err != nil {
			return handler.Result{
				Status: http.StatusInternalServerError,
				Body:   handler.Error{Message: err.Error()},
			}
		}

		if cached.Response != "" {
			ctx.Header("Content-Type", "application/json")
			ctx.Header("Cached-Response", "true")
			_, _ = ctx.Writer.Write([]byte(cached.Response))
		}
		return handler.Result{
			Status: cached.Status,
		}

	}

	result := h.F(ctx)

	marshal, err := json.Marshal(result.Body)
	if err != nil {
		sys.S.Log.Error("error to parse cache response:", err)
	}

	if err := idempotency.Update(ctx, idempotency.Idempotency{
		Id:       id,
		Endpoint: h.Endpoint,
		Status:   result.Status,
		Response: string(marshal),
	}); err != nil {
		sys.S.Log.Error("error to parse cache response:", err)
	}

	return result
}
