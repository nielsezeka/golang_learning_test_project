package quiz

import (
	"fmt"
	"learn_phase_2_local_server/db"
	"learn_phase_2_local_server/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// UpdateQuiz godoc
//
// @Summary      Update a quiz
// @Description  Update quiz fields (question, options, answers) by ID. Only provided fields will be updated. Types must match the struct: question (string), options/answers ([]string).
// @Tags         quiz
// @Accept       json
// @Produce      json
// @Param        id    path  int  true  "Quiz ID"
// @Param        quiz  body  QuizUpdateInput  true  "Quiz object (partial allowed)"  example({"question": "New Q", "options": ["A", "B"]})
// @Success      200  {object}  QuizUpdateSuccess
// @Failure      400  {object}  utils.APIError
// @Failure      404  {object}  utils.APIError
// @Failure      500  {object}  utils.APIError
// @Security     BearerAuth
// @Router       /api/quiz/{id} [put]
func UpdateQuiz(c *gin.Context) {
	id := c.Param("id")
	var input map[string]interface{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, utils.APIError{Error: err.Error()})
		return
	}

	allowedFields := map[string]func(interface{}) (interface{}, error){
		"question": utils.StringConverter,
		"options":  utils.StringArrayConverter,
		"answers":  utils.StringArrayConverter,
	}
	setClauses, args, argIdx, err := utils.BuildUpdateQuery(input, allowedFields, 1)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.APIError{Error: err.Error()})
		return
	}
	if len(setClauses) == 0 {
		c.JSON(http.StatusBadRequest, utils.APIError{Error: "no fields to update"})
		return
	}

	query := fmt.Sprintf("UPDATE quiz_table SET %s WHERE id = $%d",
		strings.Join(setClauses, ", "), argIdx)
	args = append(args, id)

	result, err := db.DB.Exec(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.APIError{Error: err.Error()})
		return
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, utils.APIError{Error: "quiz not found"})
		return
	}
	c.JSON(http.StatusOK, QuizUpdateSuccess{Message: "Quiz updated successfully"})
}

// QuizUpdateSuccess represents a successful quiz update response
// swagger:model
type QuizUpdateSuccess struct {
	Message string `json:"message"`
}
