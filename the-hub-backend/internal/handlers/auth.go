package handlers

import (
	"net/http"
	"time"

	"github.com/TheoMKgosi/The-hub/internal/config"
	"github.com/TheoMKgosi/The-hub/internal/models"
	"github.com/TheoMKgosi/The-hub/internal/util"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// ForgotPasswordRequest represents the request body for forgot password
type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email" example:"user@example.com"`
}

// ForgotPasswordResponse represents the response body for forgot password
type ForgotPasswordResponse struct {
	Message string `json:"message" example:"If an account with that email exists, a password reset link has been sent."`
}

// ResetPasswordRequest represents the request body for reset password
type ResetPasswordRequest struct {
	Token    string `json:"token" binding:"required" example:"abc123def456"`
	Password string `json:"password" binding:"required" example:"newpassword123"`
}

// ResetPasswordResponse represents the response body for reset password
type ResetPasswordResponse struct {
	Message string `json:"message" example:"Password reset successfully"`
}

// RequestPasswordReset godoc
// @Summary      Request password reset
// @Description  Send password reset email to user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body      ForgotPasswordRequest  true  "Email address"
// @Success      200      {object}  ForgotPasswordResponse
// @Failure      400      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /auth/forgot-password [post]
func RequestPasswordReset(c *gin.Context) {
	var input ForgotPasswordRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		config.Logger.Warnf("Invalid forgot password input: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	config.Logger.Infof("Password reset requested for email: %s", input.Email)

	// Check if user exists (but don't reveal if they don't)
	var user models.User
	if err := config.GetDB().Where("email = ?", input.Email).First(&user).Error; err != nil {
		// User doesn't exist, but return success to prevent email enumeration
		config.Logger.Warnf("Password reset requested for non-existent email: %s", input.Email)
		c.JSON(http.StatusOK, ForgotPasswordResponse{
			Message: "If an account with that email exists, a password reset link has been sent.",
		})
		return
	}

	// Generate reset token
	resetToken, err := util.GenerateResetToken()
	if err != nil {
		config.Logger.Errorf("Failed to generate reset token for user %s: %v", user.ID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate reset token"})
		return
	}

	// Create password reset token record
	expiresAt := time.Now().Add(1 * time.Hour) // Token expires in 1 hour
	passwordResetToken := models.PasswordResetToken{
		UserID:    user.ID,
		Token:     resetToken,
		ExpiresAt: expiresAt,
		Used:      false,
	}

	if err := config.GetDB().Create(&passwordResetToken).Error; err != nil {
		config.Logger.Errorf("Failed to create password reset token for user %s: %v", user.ID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create reset token"})
		return
	}

	// Send password reset email
	emailService := util.NewEmailService()
	if err := emailService.SendPasswordResetEmail(user.Email, resetToken); err != nil {
		config.Logger.Errorf("Failed to send password reset email to %s: %v", user.Email, err)
		// Don't return error to user as this would reveal if email exists
	}

	config.Logger.Infof("Password reset token created and email sent for user: %s", user.ID)
	c.JSON(http.StatusOK, ForgotPasswordResponse{
		Message: "If an account with that email exists, a password reset link has been sent.",
	})
}

// ResetPassword godoc
// @Summary      Reset password
// @Description  Reset user password using reset token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body      ResetPasswordRequest  true  "Reset token and new password"
// @Success      200      {object}  ResetPasswordResponse
// @Failure      400      {object}  map[string]string
// @Failure      401      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /auth/reset-password [post]
func ResetPassword(c *gin.Context) {
	var input ResetPasswordRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		config.Logger.Warnf("Invalid reset password input: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	config.Logger.Infof("Password reset attempt with token: %s", input.Token[:8]+"...")

	// Find the reset token
	var resetToken models.PasswordResetToken
	if err := config.GetDB().Where("token = ?", input.Token).First(&resetToken).Error; err != nil {
		config.Logger.Warnf("Invalid reset token used: %s", input.Token[:8]+"...")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired reset token"})
		return
	}

	// Check if token is valid
	if !resetToken.IsValid() {
		if resetToken.Used {
			config.Logger.Warnf("Attempted to use already used reset token: %s", input.Token[:8]+"...")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "This reset token has already been used"})
			return
		}
		if resetToken.IsExpired() {
			config.Logger.Warnf("Attempted to use expired reset token: %s", input.Token[:8]+"...")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "This reset token has expired"})
			return
		}
	}

	// Hash the new password
	hashedPassword, err := util.HashPassword(input.Password)
	if err != nil {
		config.Logger.Errorf("Password hashing failed for user %s: %v", resetToken.UserID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process new password"})
		return
	}

	// Update user password
	if err := config.GetDB().Model(&models.User{}).Where("id = ?", resetToken.UserID).Update("password", hashedPassword).Error; err != nil {
		config.Logger.Errorf("Failed to update password for user %s: %v", resetToken.UserID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
		return
	}

	// Mark token as used
	if err := config.GetDB().Model(&resetToken).Update("used", true).Error; err != nil {
		config.Logger.Errorf("Failed to mark reset token as used: %v", err)
		// Don't return error as password was already updated
	}

	config.Logger.Infof("Password reset successfully for user: %s", resetToken.UserID)
	c.JSON(http.StatusOK, ResetPasswordResponse{
		Message: "Password reset successfully",
	})
}

// RefreshTokenRequest represents the request body for token refresh
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required" example:"abc123def456..."`
}

// RefreshTokenResponse represents the response body for successful token refresh
type RefreshTokenResponse struct {
	Message     string `json:"message" example:"Token refreshed successfully"`
	AccessToken string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	TokenType   string `json:"token_type" example:"Bearer"`
	ExpiresIn   int    `json:"expires_in" example:"900"`
}

// RefreshToken godoc
// @Summary      Refresh access token
// @Description  Exchange refresh token for a new access token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body      RefreshTokenRequest  true  "Refresh token"
// @Success      200      {object}  RefreshTokenResponse
// @Failure      400      {object}  map[string]string
// @Failure      401      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /auth/refresh [post]
func RefreshToken(c *gin.Context) {
	var input RefreshTokenRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		config.Logger.Warnf("Invalid refresh token input: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	config.Logger.Infof("Token refresh attempt with token: %s", input.RefreshToken[:8]+"...")

	// Validate the refresh token
	refreshTokenRecord, err := util.ValidateRefreshToken(input.RefreshToken)
	if err != nil {
		if err == jwt.ErrTokenExpired {
			config.Logger.Warnf("Attempted to use expired refresh token: %s", input.RefreshToken[:8]+"...")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token has expired"})
			return
		}
		config.Logger.Warnf("Invalid refresh token used: %s", input.RefreshToken[:8]+"...")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	// Generate new access token
	accessToken, err := util.GenerateAccessToken(refreshTokenRecord.UserID)
	if err != nil {
		config.Logger.Errorf("Access token generation failed for user ID %s: %v", refreshTokenRecord.UserID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
		return
	}

	config.Logger.Infof("Token refresh successful for user: %s", refreshTokenRecord.UserID)
	c.JSON(http.StatusOK, RefreshTokenResponse{
		Message:     "Token refreshed successfully",
		AccessToken: accessToken,
		TokenType:   "Bearer",
		ExpiresIn:   900, // 15 minutes in seconds
	})
}

// LogoutRequest represents the request body for logout
type LogoutRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required" example:"abc123def456..."`
}

// LogoutResponse represents the response body for successful logout
type LogoutResponse struct {
	Message string `json:"message" example:"Logged out successfully"`
}

// Logout godoc
// @Summary      Logout user
// @Description  Revoke refresh token to log out user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body      LogoutRequest  true  "Refresh token to revoke"
// @Success      200      {object}  LogoutResponse
// @Failure      400      {object}  map[string]string
// @Failure      401      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /auth/logout [post]
func Logout(c *gin.Context) {
	var input LogoutRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		config.Logger.Warnf("Invalid logout input: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	config.Logger.Infof("Logout attempt with token: %s", input.RefreshToken[:8]+"...")

	// Validate the refresh token
	refreshTokenRecord, err := util.ValidateRefreshToken(input.RefreshToken)
	if err != nil {
		// Even if token is invalid, we return success for security
		config.Logger.Warnf("Logout attempt with invalid token: %s", input.RefreshToken[:8]+"...")
		c.JSON(http.StatusOK, LogoutResponse{
			Message: "Logged out successfully",
		})
		return
	}

	// Revoke the refresh token
	if err := util.RevokeRefreshToken(refreshTokenRecord.ID); err != nil {
		config.Logger.Errorf("Failed to revoke refresh token %s: %v", refreshTokenRecord.ID, err)
		// Don't return error to user as the token is effectively invalidated
	}

	config.Logger.Infof("User logged out successfully: %s", refreshTokenRecord.UserID)
	c.JSON(http.StatusOK, LogoutResponse{
		Message: "Logged out successfully",
	})
}
