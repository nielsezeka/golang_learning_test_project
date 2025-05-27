package quiz

import (
	"learn_phase_2_local_server/db"
	"learn_phase_2_local_server/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// DeleteQuiz godoc
//
// @Summary      Delete a quiz
// @Description  Delete a quiz by ID
// @Tags         quiz
// @Produce      json
// @Param        id  path  int  true  "Quiz ID"
// @Success      200  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Security     BearerAuth
// @Router       /api/quiz/{id} [delete]
func DeleteQuiz(c *gin.Context) {
	id := c.Param("id")
	result, err := db.DB.Exec("DELETE FROM quiz_table WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.APIError{Error: err.Error()})
		return
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, utils.APIError{Error: "quiz not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Quiz deleted successfully"})
}
