package person

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/ribgsilva/person-api/business/v1/person"
	"github.com/ribgsilva/person-api/platform/web/handler"
	"net/http"
)

// GetFilter godoc
// @Summary Search for persons
// @Description Search for persons using its fields
// @Tags Person
// @Produce json
// @Param name query person.SearchRequest false "Params"
// @Success 200 {array} person.Person
// @Failure 400 {array} handler.Error
// @Failure 404 {object} handler.Error
// @Router /v1/persons [get]
func GetFilter(ctx *gin.Context) handler.Result {

	var req person.SearchRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
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

	get, err := person.Search(ctx, req)

	switch {
	case err != nil:
		return handler.Result{
			Status: http.StatusInternalServerError,
			Body:   handler.Error{Message: err.Error()},
		}
	default:
		return handler.Result{
			Status: http.StatusOK,
			Body:   get,
		}
	}
}
