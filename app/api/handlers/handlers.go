package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/ribgsilva/person-api/app/api/handlers/v1/healthcheck"
	"github.com/ribgsilva/person-api/app/api/handlers/v1/person"
	"github.com/ribgsilva/person-api/business/web/idempotency"
	"github.com/ribgsilva/person-api/platform/web/handler"
)

func MapHandlers(r *gin.Engine) {
	g := r.Group("/v1/persons")
	{
		g.POST("", handler.Wrapper(idempotency.Handler{
			Endpoint: "create-person",
			F:        person.Post,
		}.Handle))
		g.GET("/:id", handler.Wrapper(person.Get))
		g.GET("", handler.Wrapper(person.GetFilter))
		g.PUT("/:id", handler.Wrapper(idempotency.Handler{
			Endpoint: "replace-person",
			F:        person.Put,
		}.Handle))
		g.DELETE("/:id", handler.Wrapper(idempotency.Handler{
			Endpoint: "delete-person",
			F:        person.Delete,
		}.Handle))
	}
	r.GET("/v1/healthcheck", handler.Wrapper(healthcheck.Get))
}
