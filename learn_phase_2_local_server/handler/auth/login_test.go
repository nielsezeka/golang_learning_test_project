package auth

import (
	"learn_phase_2_local_server/utils"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// Set test environment variables before anything else
	os.Setenv("JWT_SECRET", "test-jwt-secret")
	os.Setenv("REFRESH_SECRET", "test-refresh-secret")
	os.Setenv("HASH_PASS_KEY", "test-hash-key")

	// Reinitialize global variables with test values
	jwtSecret = []byte(os.Getenv("JWT_SECRET"))
	refreshSecret = []byte(os.Getenv("REFRESH_SECRET"))
	hashPassKey = os.Getenv("HASH_PASS_KEY")

	// Run tests
	code := m.Run()
	os.Exit(code)
}

func TestLogin_Success(t *testing.T) {
	ts := utils.SetupTest(t)
	defer ts.Cleanup()

	// Setup expectations
	ts.ExpectUserFound("testuser", "testpassword123", hashPassKey)

	// Make request
	ts.MakeLoginRequest(map[string]string{
		"username": "testuser",
		"password": "testpassword123",
	})

	// Call handler
	Login(ts.Context)

	// Assert response
	ts.AssertSuccessResponse(t)
}

func TestLogin_UserNotFound(t *testing.T) {
	ts := utils.SetupTest(t)
	defer ts.Cleanup()

	// Setup expectations
	ts.ExpectUserNotFound("nonexistentuser")

	// Make request
	ts.MakeLoginRequest(map[string]string{
		"username": "nonexistentuser",
		"password": "anypassword",
	})

	// Call handler
	Login(ts.Context)

	// Assert response
	ts.AssertErrorResponse(t, "User does not existed")
}

func TestLogin_InvalidPassword(t *testing.T) {
	ts := utils.SetupTest(t)
	defer ts.Cleanup()

	// Setup expectations - user exists with different password
	ts.ExpectUserFound("testuser", "correctpassword", hashPassKey)

	// Make request with wrong password
	ts.MakeLoginRequest(map[string]string{
		"username": "testuser",
		"password": "wrongpassword",
	})

	// Call handler
	Login(ts.Context)

	// Assert response
	ts.AssertErrorResponse(t, "Invalid password")
}

func TestLogin_InvalidRequestBody(t *testing.T) {
	ts := utils.SetupTest(t)
	defer ts.Cleanup()

	// Make request with invalid JSON
	ts.MakeInvalidJSONRequest()

	// Call handler
	Login(ts.Context)

	// Assert response
	ts.AssertBadRequestResponse(t, "Invalid request")
}

func TestLogin_MissingFields(t *testing.T) {
	testCases := []struct {
		name string
		body map[string]string
	}{
		{
			name: "Missing username",
			body: map[string]string{"password": "testpass"},
		},
		{
			name: "Missing password",
			body: map[string]string{"username": "testuser"},
		},
		{
			name: "Empty username",
			body: map[string]string{"username": "", "password": "testpass"},
		},
		{
			name: "Empty password",
			body: map[string]string{"username": "testuser", "password": ""},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ts := utils.SetupTest(t)
			defer ts.Cleanup()
			// Setup expectations - user lookup will be called even with empty fields
			username := tc.body["username"]
			ts.ExpectUserNotFound(username)
			// Make request
			ts.MakeLoginRequest(tc.body)
			// Call handler
			Login(ts.Context)
			// Assert response
			ts.AssertErrorResponse(t, "User does not existed")
		})
	}
}
