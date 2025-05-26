package router_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"learn_phase_2_local_server/router"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRootEndpoint(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := router.SetupRouter()

	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "Hello, Gin!")
}

func TestProtectedQuizRoute_Unauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := router.SetupRouter()

	req, _ := http.NewRequest("GET", "/api/quiz", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, 401, w.Code)
	assert.Contains(t, w.Body.String(), "error")
}
func TestProtectedQuizRoute_LoginSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := router.SetupRouter()
	loginReqBody := `{"username":"admin","password":"password123"}`
	req, _ := http.NewRequest("POST", "/api/login",
		strings.NewReader(loginReqBody))
	req.Header.Set("Content-Type", "application/json")
	loginW := httptest.NewRecorder()
	r.ServeHTTP(loginW, req)

	assert.Equal(t, 200, loginW.Code)
	var resp map[string]string
	err := json.Unmarshal(loginW.Body.Bytes(), &resp)
	assert.NoError(t, err)
	_, ok := resp["token"]
	assert.True(t, ok)
}
func TestProtectedQuizRoute_LoginFail(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := router.SetupRouter()
	loginReqBody := `{"username":"admin","password":"wrongpass"}`
	req, _ := http.NewRequest("POST", "/api/login",
		strings.NewReader(loginReqBody))
	req.Header.Set("Content-Type", "application/json")
	loginW := httptest.NewRecorder()
	r.ServeHTTP(loginW, req)
	assert.Equal(t, 401, loginW.Code)
}
