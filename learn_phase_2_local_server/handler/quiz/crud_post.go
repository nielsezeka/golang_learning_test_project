package quiz

import (
	"learn_phase_2_local_server/db"
	"learn_phase_2_local_server/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

// PostQuiz godoc
//
// @Summary      Create a new quiz
// @Description  Create a new quiz with question, options, and answers (all required, options/answers must be arrays of strings)
// @Tags         quiz
// @Accept       json
// @Produce      json
// @Param        quiz  body  object  true  "Quiz object"  example({"question": "What is the capital?", "options": ["A", "B"], "answers": ["A"]})
// @Success      201  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Security     BearerAuth
// @Router       /api/quiz [post]
func PostQuiz(c *gin.Context) {
	var quiz struct {
		Question string   `json:"question"`
		Options  []string `json:"options"`
		Answers  []string `json:"answers"`
	}

	if err := c.ShouldBindJSON(&quiz); err != nil {
		c.JSON(http.StatusBadRequest, utils.APIError{Error: err.Error()})
		return
	}
	var id int
	err := db.DB.QueryRow(
		`INSERT INTO quiz_table (question, options, answers) 
		 VALUES ($1, $2, $3) RETURNING id`,
		quiz.Question, pq.Array(quiz.Options), pq.Array(quiz.Answers),
	).Scan(&id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.APIError{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Quiz created successfully",
		"quiz": gin.H{
			"id":       id,
			"question": quiz.Question,
			"options":  quiz.Options,
			"answers":  quiz.Answers,
		},
	})
}
