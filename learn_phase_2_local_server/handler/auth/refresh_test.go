package auth

import (
	"learn_phase_2_local_server/utils"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestRefresh_Success(t *testing.T) {
	ts := utils.SetupTest(t, &refreshTokens)
	defer ts.Cleanup()

	userID := 123
	refreshToken := ts.CreateValidRefreshToken(userID, refreshSecret, &refreshTokens)

	// Make request
	ts.MakeRefreshRequest(refreshToken)

	// Call handler
	Refresh(ts.Context)

	// Assert response
	ts.AssertRefreshSuccessResponse(t)
}

func TestRefresh_InvalidRefreshToken(t *testing.T) {
	ts := utils.SetupTest(t, &refreshTokens)
	defer ts.Cleanup()

	// Make request with invalid token
	ts.MakeRefreshRequest("invalid_token")

	// Call handler
	Refresh(ts.Context)

	// Assert response
	ts.AssertUnauthorizedResponse(t, "Invalid refresh token")
}

func TestRefresh_ExpiredRefreshToken(t *testing.T) {
	ts := utils.SetupTest(t, &refreshTokens)
	defer ts.Cleanup()

	userID := 123
	expiredToken := ts.CreateExpiredRefreshToken(userID, refreshSecret)

	// Make request with expired token
	ts.MakeRefreshRequest(expiredToken)

	// Call handler
	Refresh(ts.Context)

	// Assert response
	ts.AssertUnauthorizedResponse(t, "Invalid refresh token")
}

func TestRefresh_TokenNotRecognized(t *testing.T) {
	ts := utils.SetupTest(t, &refreshTokens)
	defer ts.Cleanup()

	userID := 123
	// Create a valid token but don't store it in refreshTokens map
	validToken := ts.CreateValidRefreshToken(userID, refreshSecret, nil) // Don't store in map

	// Make request with unrecognized token
	ts.MakeRefreshRequest(validToken)

	// Call handler
	Refresh(ts.Context)

	// Assert response
	ts.AssertUnauthorizedResponse(t, "Refresh token not recognized")
}

func TestRefresh_InvalidRequestBody(t *testing.T) {
	ts := utils.SetupTest(t, &refreshTokens)
	defer ts.Cleanup()

	// Make request with invalid JSON
	ts.MakeInvalidJSONRequest()

	// Call handler
	Refresh(ts.Context)

	// Assert response
	ts.AssertBadRequestResponse(t, "Invalid request")
}

func TestRefresh_MissingRefreshToken(t *testing.T) {
	ts := utils.SetupTest(t, &refreshTokens)
	defer ts.Cleanup()

	// Make request with empty refresh token
	ts.MakeRefreshRequest("")

	// Call handler
	Refresh(ts.Context)

	// Assert response
	ts.AssertUnauthorizedResponse(t, "Invalid refresh token")
}

func TestRefresh_TokenWithInvalidSignature(t *testing.T) {
	ts := utils.SetupTest(t, &refreshTokens)
	defer ts.Cleanup()

	userID := 123
	// Create token with wrong secret
	wrongToken := ts.CreateExpiredRefreshToken(userID, []byte("wrong-secret"))

	// Make request with token signed with wrong secret
	ts.MakeRefreshRequest(wrongToken)

	// Call handler
	Refresh(ts.Context)

	// Assert response
	ts.AssertUnauthorizedResponse(t, "Invalid refresh token")
}

func TestRefresh_TokenWithoutUserID(t *testing.T) {
	ts := utils.SetupTest(t, &refreshTokens)
	defer ts.Cleanup()

	// Create a token without user_id claim
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Add(time.Hour * 24 * 7).Unix(),
		// Missing user_id claim
	})
	tokenString, _ := token.SignedString(refreshSecret)
	refreshTokens[tokenString] = 123 // Store in map

	// Make request
	ts.MakeRefreshRequest(tokenString)

	// Call handler
	Refresh(ts.Context)

	// Assert response
	ts.AssertUnauthorizedResponse(t, "Invalid refresh token claims")
}
