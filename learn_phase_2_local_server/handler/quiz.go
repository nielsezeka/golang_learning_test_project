package handler

import (
	"fmt"
	"learn_phase_2_local_server/db"
	"learn_phase_2_local_server/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

// GetQuiz godoc
//
//	@Summary		Get all quizzes
//	@Description	Returns a list of all quizzes
//	@Tags			quiz
//	@Produce		json
//	@Success		200	{array}		map[string]interface{}
//	@Failure		500	{object}	map[string]string
//	@Security		BearerAuth
//	@Router			/api/quiz [get]
func GetQuiz(c *gin.Context) {
	rows, _ := db.DB.Query("SELECT * FROM quiz_table")

	var id int
	var question string
	var options, answers []string
	var quizzes []map[string]interface{}
	for rows.Next() {
		err := rows.Scan(&id, &question, pq.Array(&options), pq.Array(&answers))
		if err != nil {
			utils.ErrorReturnHandler(c, http.StatusInternalServerError, err)
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
		utils.ErrorReturnHandler(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, quizzes)
}

// PostQuiz godoc
//
//	@Summary		Create a new quiz
//	@Description	Create a new quiz with question, options, and answers (all required, options/answers must be arrays of strings)
//	@Tags			quiz
//	@Accept			json
//	@Produce		json
//	@Param			quiz	body		object	true	"Quiz object"	example({"question": "What is the capital?", "options": ["A", "B"], "answers": ["A"]})
//	@Success		201		{object}	map[string]interface{}
//	@Failure		400		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Security		BearerAuth
//	@Router			/api/quiz [post]
func PostQuiz(c *gin.Context) {
	var quiz struct {
		Question string   `json:"question"`
		Options  []string `json:"options"`
		Answers  []string `json:"answers"`
	}

	if err := c.ShouldBindJSON(&quiz); err != nil {
		utils.ErrorReturnHandler(c, http.StatusBadRequest, err)
		return
	}
	var id int
	err := db.DB.QueryRow(
		`INSERT INTO quiz_table (question, options, answers) 
		 VALUES ($1, $2, $3) RETURNING id`,
		quiz.Question, pq.Array(quiz.Options), pq.Array(quiz.Answers),
	).Scan(&id)
	if err != nil {
		utils.ErrorReturnHandler(c, http.StatusInternalServerError, err)
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

// DeleteQuiz godoc
//
//	@Summary		Delete a quiz
//	@Description	Delete a quiz by ID
//	@Tags			quiz
//	@Produce		json
//	@Param			id	path		int	true	"Quiz ID"
//	@Success		200	{object}	map[string]string
//	@Failure		404	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Security		BearerAuth
//	@Router			/api/quiz/{id} [delete]
func DeleteQuiz(c *gin.Context) {
	id := c.Param("id")
	result, err := db.DB.Exec("DELETE FROM quiz_table WHERE id = $1", id)
	if err != nil {
		utils.ErrorReturnHandler(c, http.StatusInternalServerError, err)
		return
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		utils.ErrorReturnHandler(c, http.StatusNotFound, fmt.Errorf("quiz not found"))
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Quiz deleted successfully"})
}

// UpdateQuiz godoc
//
//	@Summary		Update a quiz
//	@Description	Update quiz fields (question, options, answers) by ID. Only provided fields will be updated. Types must match the struct: question (string), options/answers ([]string).
//	@Tags			quiz
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int		true	"Quiz ID"
//	@Param			quiz	body		object	true	"Quiz object (partial allowed)"	example({"question": "New Q", "options": ["A", "B"]})
//	@Success		200		{object}	map[string]string
//	@Failure		400		{object}	map[string]string
//	@Failure		404		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Security		BearerAuth
//	@Router			/api/quiz/{id} [put]
func UpdateQuiz(c *gin.Context) {
	id := c.Param("id")
	var input map[string]interface{}
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorReturnHandler(c, http.StatusBadRequest, err)
		return
	}

	allowedFields := map[string]func(interface{}) (interface{}, error){
		"question": utils.StringConverter,
		"options":  utils.StringArrayConverter,
		"answers":  utils.StringArrayConverter,
	}
	setClauses, args, argIdx, err := utils.BuildUpdateQuery(input, allowedFields, 1)
	if err != nil {
		utils.ErrorReturnHandler(c, http.StatusBadRequest, err)
		return
	}
	if len(setClauses) == 0 {
		utils.ErrorReturnHandler(c, http.StatusBadRequest, fmt.Errorf("no fields to update"))
		return
	}

	query := fmt.Sprintf("UPDATE quiz_table SET %s WHERE id = $%d",
		strings.Join(setClauses, ", "), argIdx)
	args = append(args, id)

	result, err := db.DB.Exec(query, args...)
	if err != nil {
		utils.ErrorReturnHandler(c, http.StatusInternalServerError, err)
		return
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		utils.ErrorReturnHandler(c, http.StatusNotFound, fmt.Errorf("quiz not found"))
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Quiz updated successfully"})
}
