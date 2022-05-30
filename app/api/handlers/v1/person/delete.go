package person

import (
	"github.com/gin-gonic/gin"
	"github.com/ribgsilva/person-api/business/v1/person"
	"github.com/ribgsilva/person-api/platform/web/handler"
	"net/http"
)

// Delete godoc
// @Summary Delete a person
// @Description Delete a person using its id
// @Tags Person
// @Produce json
// @Param id path string true "Person id"
// @Success 200 {object} person.Person
// @Failure 400 {array} handler.Error
// @Failure 404 {object} handler.Error
// @Router /v1/persons/{id} [delete]
func Delete(ctx *gin.Context) handler.Result {

	delete, err := person.Remove(ctx, ctx.Param("id"))

	switch {
	case err != nil:
		return handler.Result{
			Status: http.StatusInternalServerError,
			Body:   handler.Error{Message: err.Error()},
		}
	case delete.Id == "":
		return handler.Result{
			Status: http.StatusNotFound,
			Body:   handler.Error{Message: "person not found"},
		}
	default:
		return handler.Result{
			Status: http.StatusOK,
			Body:   delete,
		}
	}
}
