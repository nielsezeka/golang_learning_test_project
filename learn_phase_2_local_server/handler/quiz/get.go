package quiz

import (
	"learn_phase_2_local_server/db"
	"learn_phase_2_local_server/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

// GetQuiz godoc
//
// @Summary      Get all quizzes
// @Description  Returns a list of all quizzes
// @Tags         quiz
// @Produce      json
// @Success      200  {array}  map[string]interface{}
// @Failure      500  {object}  map[string]string
// @Security     BearerAuth
// @Router       /api/quiz [get]
func GetQuiz(c *gin.Context) {
	rows, _ := db.DB.Query("SELECT * FROM quiz_table")

	var id int
	var question string
	var options, answers []string
	var quizzes []map[string]interface{}
	for rows.Next() {
		err := rows.Scan(&id, &question, pq.Array(&options), pq.Array(&answers))
		if err != nil {
			c.JSON(http.StatusInternalServerError, utils.APIError{Error: err.Error()})
			return
		}
		quiz := map[string]interface{}{
			"id":       id,
			"question": question,
			"options":  options,
			"answers":  answers,
		}
		quizzes = append(quizzes, quiz)
	}
	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, utils.APIError{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, quizzes)
}
