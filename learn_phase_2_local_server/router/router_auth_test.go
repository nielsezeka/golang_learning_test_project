package router_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"learn_phase_2_local_server/router"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestProtectedQuizRoute_Unauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := router.SetupRouter()

	req, _ := http.NewRequest("GET", "/api/quiz/", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	fmt.Println("Request URL:", req.URL)
	assert.Equal(t, 401, w.Code)
	assert.Contains(t, w.Body.String(), "error")
}
