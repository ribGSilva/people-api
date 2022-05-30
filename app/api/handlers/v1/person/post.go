package person

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/ribgsilva/person-api/business/v1/person"
	"github.com/ribgsilva/person-api/platform/web/handler"
	"net/http"
)

// Post godoc
// @Summary Creates a new "person"
// @Description Create a new "person" into the system
// @Tags Person
// @Produce json
// @Param X-Idempotency-Key header string true "Idempotency"
// @Param Request body person.CreateRequest true "Request body"
// @Success 201 {object} person.CreateResponse
// @Failure 400 {array} handler.Error
// @Router /v1/persons [post]
func Post(ctx *gin.Context) handler.Result {

	var req person.CreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			out := make([]handler.Error, len(ve))
			for i, fe := range ve {
				out[i] = handler.Error{Field: fe.Field(), Message: fe.Error()}
			}
			return handler.Result{
				Status: http.StatusBadRequest,
				Body:   out,
			}
		} else {
			return handler.Result{
				Status: http.StatusInternalServerError,
				Body:   handler.Error{Message: err.Error()},
			}
		}
	}

	if req.Type == "employee" && req.Role == nil {
		return handler.Result{
			Status: http.StatusBadRequest,
			Body: []handler.Error{
				{
					Message: "employees must have a role",
				},
			},
		}
	}

	if req.Type == "contractor" && req.ContactDuration == nil {
		return handler.Result{
			Status: http.StatusBadRequest,
			Body: []handler.Error{
				{
					Message: "contractor must have a contract duration",
				},
			},
		}
	}

	create, err := person.Create(ctx, req)
	switch {
	case err != nil:
		return handler.Result{
			Status: http.StatusInternalServerError,
			Body:   handler.Error{Message: err.Error()},
		}
	default:
		return handler.Result{
			Status: http.StatusCreated,
			Body:   create,
		}
	}
}
