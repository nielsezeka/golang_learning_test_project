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
