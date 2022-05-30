package person

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/ribgsilva/person-api/business/v1/person"
	"github.com/ribgsilva/person-api/platform/web/handler"
	"net/http"
)

// Put godoc
// @Summary Replaces the person
// @Description Replaces all data of a person
// @Tags Person
// @Produce json
// @Param X-Idempotency-Key header string true "Idempotency"
// @Param id path string true "Person id"
// @Param Request body person.UpdateRequest true "Request body"
// @Success 200 {object} person.Person
// @Failure 400 {array} handler.Error
// @Failure 404 {array} handler.Error
// @Router /v1/persons/{id} [put]
func Put(ctx *gin.Context) handler.Result {

	var req person.UpdateRequest
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

	replace, err := person.Replace(ctx, ctx.Param("id"), req)
	switch {
	case err != nil:
		return handler.Result{
			Status: http.StatusInternalServerError,
			Body:   handler.Error{Message: err.Error()},
		}
	case replace.Id == "":
		return handler.Result{
			Status: http.StatusNotFound,
			Body:   handler.Error{Message: "person not found"},
		}
	default:
		return handler.Result{
			Status: http.StatusOK,
			Body:   replace,
		}
	}
}
