package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// Define an interface for DB operations
type QuizDB interface {
	GetQuiz(id int) (string, error)
}

// Mock implementation of QuizDB
type MockQuizDB struct {
	QuizData map[int]string
}

func (m *MockQuizDB) GetQuiz(id int) (string, error) {
	if quiz, ok := m.QuizData[id]; ok {
		return quiz, nil
	}
	return "", nil
}

// Example handler function that uses the interface
func GetQuizHandler(db QuizDB, id int) (string, error) {
	return db.GetQuiz(id)
}

func TestQuizAPI(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	api := r.Group("/api")
	{
		api.GET("/quiz", GetQuiz)
		api.POST("/quiz", PostQuiz)
	}

	// Test GET /api/quiz
	req, _ := http.NewRequest("GET", "/api/quiz", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}
