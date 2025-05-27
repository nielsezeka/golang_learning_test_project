package utils

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"learn_phase_2_local_server/db"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

// TestSetup contains all the mock objects and helper methods for testing
type TestSetup struct {
	MockDB   sqlmock.Sqlmock
	Recorder *httptest.ResponseRecorder
	Context  *gin.Context
	cleanup  func()
}

// SetupTest initializes a test environment with mock database and Gin context
// Supports optional refreshTokens map for auth testing
func SetupTest(t *testing.T, refreshTokens ...*map[string]int) *TestSetup {
	// Create mock DB
	mockDBConn, mock, err := sqlmock.New()
	assert.NoError(t, err)

	// Replace the global DB with mock
	originalDB := db.DB
	db.DB = mockDBConn

	// Handle refresh tokens isolation if provided
	var originalRefreshTokens map[string]int
	var refreshTokensMap *map[string]int
	if len(refreshTokens) > 0 && refreshTokens[0] != nil {
		refreshTokensMap = refreshTokens[0]
		// Store original refreshTokens and clear it for test isolation
		originalRefreshTokens = make(map[string]int)
		for k, v := range *refreshTokensMap {
			originalRefreshTokens[k] = v
		}
		*refreshTokensMap = make(map[string]int)
	}

	// Setup Gin test context
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	return &TestSetup{
		MockDB:   mock,
		Recorder: w,
		Context:  c,
		cleanup: func() {
			mockDBConn.Close()
			db.DB = originalDB
			// Restore original refreshTokens if provided
			if refreshTokensMap != nil {
				*refreshTokensMap = originalRefreshTokens
			}
		},
	}
}

// Cleanup properly cleans up test resources
func (ts *TestSetup) Cleanup() {
	ts.cleanup()
}

// =============================================================================
// REQUEST HELPERS
// =============================================================================

// MakeJSONRequest creates a request with JSON body for any endpoint
func (ts *TestSetup) MakeJSONRequest(method, path string, body interface{}) {
	var buf *bytes.Buffer
	if body != nil {
		jsonBody, _ := json.Marshal(body)
		buf = bytes.NewBuffer(jsonBody)
	} else {
		buf = bytes.NewBuffer([]byte{})
	}

	ts.Context.Request, _ = http.NewRequest(method, path, buf)
	ts.Context.Request.Header.Set("Content-Type", "application/json")
}

// MakeLoginRequest creates a POST request with JSON body for login endpoint
func (ts *TestSetup) MakeLoginRequest(body map[string]string) {
	ts.MakeJSONRequest("POST", "/login", body)
}

// MakeRefreshRequest creates a POST request with JSON body for refresh endpoint
func (ts *TestSetup) MakeRefreshRequest(refreshToken string) {
	body := map[string]string{"refresh_token": refreshToken}
	ts.MakeJSONRequest("POST", "/refresh", body)
}

// MakeInvalidJSONRequest creates a request with malformed JSON for /login by default
func (ts *TestSetup) MakeInvalidJSONRequest() {
	ts.MakeInvalidJSONRequestFor("/login")
}

// MakeInvalidJSONRequestFor creates a request with malformed JSON for any endpoint
func (ts *TestSetup) MakeInvalidJSONRequestFor(endpoint string) {
	ts.Context.Request, _ = http.NewRequest("POST", endpoint, bytes.NewBuffer([]byte("invalid json")))
	ts.Context.Request.Header.Set("Content-Type", "application/json")
}

// =============================================================================
// DATABASE MOCKING HELPERS
// =============================================================================

// ExpectUserQuery sets up the expected SQL query for user lookup
func (ts *TestSetup) ExpectUserQuery(username string) *sqlmock.ExpectedQuery {
	return ts.MockDB.ExpectQuery("SELECT id, username, password FROM users WHERE username = \\$1").
		WithArgs(username)
}

// ExpectUserNotFound mocks a scenario where user is not found in database
func (ts *TestSetup) ExpectUserNotFound(username string) {
	ts.ExpectUserQuery(username).WillReturnError(sqlmock.ErrCancelled)
}

// ExpectUserFound mocks a scenario where user exists with given username and password
func (ts *TestSetup) ExpectUserFound(username, password string, hashKey ...string) {
	key := ""
	if len(hashKey) > 0 {
		key = hashKey[0]
	}
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password+key), bcrypt.DefaultCost)
	rows := sqlmock.NewRows([]string{"id", "username", "password"}).
		AddRow(1, username, string(hashedPassword))
	ts.ExpectUserQuery(username).WillReturnRows(rows)
}

// =============================================================================
// TOKEN HELPERS
// =============================================================================

// CreateValidRefreshToken creates a valid refresh token for testing and stores it in refreshTokens map
func (ts *TestSetup) CreateValidRefreshToken(userID int, secret []byte, refreshTokensMap *map[string]int) string {
	refreshToken, _ := ts.createToken(userID, secret, time.Hour*24*7)
	if refreshTokensMap != nil {
		(*refreshTokensMap)[refreshToken] = userID
	}
	return refreshToken
}

// CreateExpiredRefreshToken creates an expired refresh token for testing
func (ts *TestSetup) CreateExpiredRefreshToken(userID int, secret []byte) string {
	refreshToken, _ := ts.createToken(userID, secret, -time.Hour) // Expired 1 hour ago
	return refreshToken
}

// createToken is a helper function to create JWT tokens for testing
func (ts *TestSetup) createToken(userID int, secret []byte, duration time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(duration).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

// =============================================================================
// RESPONSE ASSERTION HELPERS
// =============================================================================

// AssertResponse validates the HTTP response and returns the parsed JSON body
func (ts *TestSetup) AssertResponse(t *testing.T, expectedStatus int) map[string]interface{} {
	assert.Equal(t, expectedStatus, ts.Recorder.Code)
	assert.NoError(t, ts.MockDB.ExpectationsWereMet())

	var response map[string]interface{}
	err := json.Unmarshal(ts.Recorder.Body.Bytes(), &response)
	assert.NoError(t, err)
	return response
}

// AssertErrorWithStatus validates an error response with a given status and error message
func (ts *TestSetup) AssertErrorWithStatus(t *testing.T, status int, expectedError string) {
	response := ts.AssertResponse(t, status)
	assert.Equal(t, expectedError, response["error"])
}

// AssertSuccessResponse validates a successful login response with tokens
func (ts *TestSetup) AssertSuccessResponse(t *testing.T) {
	response := ts.AssertResponse(t, http.StatusOK)
	assert.Contains(t, response, "token")
	assert.Contains(t, response, "refresh_token")
	assert.NotEmpty(t, response["token"])
	assert.NotEmpty(t, response["refresh_token"])
}

// AssertRefreshSuccessResponse validates a successful refresh response with new token
func (ts *TestSetup) AssertRefreshSuccessResponse(t *testing.T) {
	response := ts.AssertResponse(t, http.StatusOK)
	assert.Contains(t, response, "token")
	assert.NotEmpty(t, response["token"])
	// Refresh endpoint only returns new access token, not refresh token
	assert.NotContains(t, response, "refresh_token")
}

// AssertErrorResponse validates an unauthorized error response
func (ts *TestSetup) AssertErrorResponse(t *testing.T, expectedError string) {
	ts.AssertErrorWithStatus(t, http.StatusUnauthorized, expectedError)
}

// AssertBadRequestResponse validates a bad request error response
func (ts *TestSetup) AssertBadRequestResponse(t *testing.T, expectedError string) {
	ts.AssertErrorWithStatus(t, http.StatusBadRequest, expectedError)
}

// AssertUnauthorizedResponse validates an unauthorized error response (alias for AssertErrorResponse)
func (ts *TestSetup) AssertUnauthorizedResponse(t *testing.T, expectedError string) {
	ts.AssertErrorWithStatus(t, http.StatusUnauthorized, expectedError)
}

// AssertInternalServerErrorResponse validates a server error response
func (ts *TestSetup) AssertInternalServerErrorResponse(t *testing.T, expectedError string) {
	ts.AssertErrorWithStatus(t, http.StatusInternalServerError, expectedError)
}
