package handler

import (
	"fmt"
	"net/http"

	"learn_phase_2_local_server/db"

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
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var id int
	err := db.DB.QueryRow(
		`INSERT INTO quiz_table (question, options, answers) 
		 VALUES ($1, $2, $3) RETURNING id`,
		quiz.Question, pq.Array(quiz.Options), pq.Array(quiz.Answers),
	).Scan(&id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Quiz not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Quiz deleted successfully"})
}

func UpdateQuiz(c *gin.Context) {
	id := c.Param("id")
	var input map[string]interface{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	setClauses := []string{}
	args := []interface{}{}
	argIdx := 1

	if question, ok := input["question"].(string); ok {
		setClauses = append(setClauses, fmt.Sprintf("question = $%d", argIdx))
		args = append(args, question)
		argIdx++
	}
	if options, ok := input["options"].([]interface{}); ok {
		strOptions := make([]string, len(options))
		for i, v := range options {
			strOptions[i], _ = v.(string)
		}
		setClauses = append(setClauses, fmt.Sprintf("options = $%d", argIdx))
		args = append(args, pq.Array(strOptions))
		argIdx++
	}
	if answers, ok := input["answers"].([]interface{}); ok {
		strAnswers := make([]string, len(answers))
		for i, v := range answers {
			strAnswers[i], _ = v.(string)
		}
		setClauses = append(setClauses, fmt.Sprintf("answers = $%d", argIdx))
		args = append(args, pq.Array(strAnswers))
		argIdx++
	}

	if len(setClauses) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No fields to update"})
		return
	}

	query := fmt.Sprintf("UPDATE quiz_table SET %s WHERE id = $%d",
		strings.Join(setClauses, ", "), argIdx)
	args = append(args, id)

	result, err := db.DB.Exec(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Quiz not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Quiz updated successfully"})
}
