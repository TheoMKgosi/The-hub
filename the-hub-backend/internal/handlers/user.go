package handlers

import (
	"net/http"
	"time"

	"github.com/TheoMKgosi/The-hub/internal/config"
	"github.com/TheoMKgosi/The-hub/internal/models"
	"github.com/TheoMKgosi/The-hub/internal/util"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GetUser godoc
// @Summary      Get user profile
// @Description  Fetch a specific user's profile
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        ID   path      int  true  "User ID"
// @Success      200  {object}  models.User
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /users/{ID} [get]
func GetUser(c *gin.Context) {
	userIDStr := c.Param("ID")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		config.Logger.Warnf("Invalid user ID param: %s", userIDStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Check if the requesting user is accessing their own profile or is admin
	requestingUserID, exist := c.Get("userID")
	if !exist {
		config.Logger.Warn("userID not found in context during user fetch")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	requestingUserIDUUID, ok := requestingUserID.(uuid.UUID)
	if !ok {
		config.Logger.Errorf("Invalid userID type in context: %T", requestingUserID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// Users can only access their own profile (unless they're admin)
	if requestingUserIDUUID != userID {
		config.Logger.Warnf("User %s attempted to access user %s profile (forbidden)", requestingUserIDUUID, userID)
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only access your own profile"})
		return
	}

	var user models.User
	if err := config.GetDB().First(&user, userID).Error; err != nil {
		config.Logger.Warnf("User not found: ID %s", userID)
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Remove password from response
	user.Password = ""
	config.Logger.Infof("User profile retrieved successfully: ID %s", userID)
	c.JSON(http.StatusOK, user)
}

// GetUserSettings godoc
// @Summary      Get user settings
// @Description  Fetch a specific user's settings
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        ID   path      int  true  "User ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /users/{ID}/settings [get]
func GetUserSettings(c *gin.Context) {
	userIDStr := c.Param("ID")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		config.Logger.Warnf("Invalid user ID param: %s", userIDStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Check if the requesting user is accessing their own settings or is admin
	requestingUserID, exist := c.Get("userID")
	if !exist {
		config.Logger.Warn("userID not found in context during settings fetch")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	requestingUserIDUUID, ok := requestingUserID.(uuid.UUID)
	if !ok {
		config.Logger.Errorf("Invalid userID type in context: %T", requestingUserID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// Users can only access their own settings (unless they're admin)
	if requestingUserIDUUID != userID {
		config.Logger.Warnf("User %s attempted to access user %s settings (forbidden)", requestingUserIDUUID, userID)
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only access your own settings"})
		return
	}

	var user models.User
	if err := config.GetDB().First(&user, userID).Error; err != nil {
		config.Logger.Warnf("User not found: ID %s", userID)
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	config.Logger.Infof("User settings retrieved successfully: ID %s", userID)
	c.JSON(http.StatusOK, gin.H{"settings": user.Settings})
}

// UpdateUserSettings godoc
// @Summary      Update user settings
// @Description  Update a specific user's settings
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        ID       path      int                      true  "User ID"
// @Param        settings body      map[string]interface{}   true  "Settings object"
// @Success      200      {object}  map[string]interface{}
// @Failure      400      {object}  map[string]string
// @Failure      401      {object}  map[string]string
// @Failure      403      {object}  map[string]string
// @Failure      404      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /users/{ID}/settings [put]
func UpdateUserSettings(c *gin.Context) {
	userIDStr := c.Param("ID")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		config.Logger.Warnf("Invalid user ID param: %s", userIDStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Check if the requesting user is updating their own settings or is admin
	requestingUserID, exist := c.Get("userID")
	if !exist {
		config.Logger.Warn("userID not found in context during settings update")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	requestingUserIDUUID, ok := requestingUserID.(uuid.UUID)
	if !ok {
		config.Logger.Errorf("Invalid userID type in context: %T", requestingUserID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// Users can only update their own settings (unless they're admin)
	if requestingUserIDUUID != userID {
		config.Logger.Warnf("User %s attempted to update user %s settings (forbidden)", requestingUserIDUUID, userID)
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only update your own settings"})
		return
	}

	var user models.User
	if err := config.GetDB().First(&user, userID).Error; err != nil {
		config.Logger.Warnf("User not found: ID %s", userID)
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var input map[string]interface{}
	if err := c.ShouldBindJSON(&input); err != nil {
		config.Logger.Warnf("Invalid settings input for user ID %d: %v", userID, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	// Update settings
	if err := config.GetDB().Model(&user).Update("settings", input).Error; err != nil {
		config.Logger.Errorf("Error updating user settings ID %d: %v", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update settings"})
		return
	}

	config.Logger.Infof("User settings updated successfully: ID %d", userID)
	c.JSON(http.StatusOK, gin.H{"settings": input})
}

// PatchUserSettings godoc
// @Summary      Partially update user settings
// @Description  Update specific user settings without replacing the entire settings object
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        ID       path      int                      true  "User ID"
// @Param        settings body      map[string]interface{}   true  "Partial settings object"
// @Success      200      {object}  map[string]interface{}
// @Failure      400      {object}  map[string]string
// @Failure      401      {object}  map[string]string
// @Failure      403      {object}  map[string]string
// @Failure      404      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /users/{ID}/settings [patch]
func PatchUserSettings(c *gin.Context) {
	userIDStr := c.Param("ID")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		config.Logger.Warnf("Invalid user ID param: %s", userIDStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Check if the requesting user is updating their own settings or is admin
	requestingUserID, exist := c.Get("userID")
	if !exist {
		config.Logger.Warn("userID not found in context during settings patch")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	requestingUserIDUUID, ok := requestingUserID.(uuid.UUID)
	if !ok {
		config.Logger.Errorf("Invalid userID type in context: %T", requestingUserID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// Users can only update their own settings (unless they're admin)
	if requestingUserIDUUID != userID {
		config.Logger.Warnf("User %s attempted to patch user %s settings (forbidden)", requestingUserIDUUID, userID)
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only update your own settings"})
		return
	}

	var user models.User
	if err := config.GetDB().First(&user, userID).Error; err != nil {
		config.Logger.Warnf("User not found: ID %s", userID)
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var input map[string]interface{}
	if err := c.ShouldBindJSON(&input); err != nil {
		config.Logger.Warnf("Invalid settings input for user ID %d: %v", userID, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	// Initialize settings if nil
	if user.Settings == nil {
		user.Settings = make(models.UserSettings)
	}

	// Merge the new settings with existing ones
	for key, value := range input {
		user.Settings[key] = value
	}

	// Update settings
	if err := config.GetDB().Model(&user).Update("settings", user.Settings).Error; err != nil {
		config.Logger.Errorf("Error patching user settings ID %d: %v", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update settings"})
		return
	}

	config.Logger.Infof("User settings patched successfully: ID %d", userID)
	c.JSON(http.StatusOK, gin.H{"settings": user.Settings})
}

// GetUsers godoc
// @Summary      Get all users
// @Description  Fetch all users (admin only)
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  map[string][]models.User
// @Failure      401  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /users [get]
func GetUsers(c *gin.Context) {
	var users []models.User

	config.Logger.Info("Fetching all users")
	result := config.GetDB().Find(&users)

	if result.Error != nil {
		config.Logger.Errorf("Error fetching users: %v", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch users"})
		return
	}

	config.Logger.Infof("Successfully fetched %d users", len(users))
	c.JSON(http.StatusOK, gin.H{"users": users})
}

// LoginRequest represents the request body for user login
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email" example:"user@example.com"`
	Password string `json:"password" binding:"required" example:"password123"`
}

// LoginResponse represents the response body for successful login
type LoginResponse struct {
	Message      string      `json:"message" example:"Login successful"`
	AccessToken  string      `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	RefreshToken string      `json:"refresh_token" example:"abc123def456..."`
	TokenType    string      `json:"token_type" example:"Bearer"`
	ExpiresIn    int         `json:"expires_in" example:"900"`
	User         models.User `json:"user"`
}

// Login godoc
// @Summary      User login
// @Description  Authenticate user and return JWT token
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user  body      LoginRequest  true  "Login credentials"
// @Success      200   {object}  LoginResponse
// @Failure      400   {object}  map[string]string
// @Failure      401   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /auth/login [post]
func Login(c *gin.Context) {
	var input LoginRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		config.Logger.Warnf("Invalid login input: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	config.Logger.Infof("Login attempt for email: %s", input.Email)
	var user models.User
	if err := config.GetDB().Where("email = ?", input.Email).First(&user).Error; err != nil {
		config.Logger.Warnf("User not found with email: %s", input.Email)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	if !util.CheckPasswordHash(user.Password, input.Password) {
		config.Logger.Warnf("Password mismatch for user: %s (ID: %d)", input.Email, user.ID)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Generate access token (short-lived)
	accessToken, err := util.GenerateAccessToken(user.ID)
	if err != nil {
		config.Logger.Errorf("Access token generation failed for user ID %s: %v", user.ID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
		return
	}

	// Generate refresh token (long-lived)
	refreshTokenString, err := util.GenerateRefreshToken()
	if err != nil {
		config.Logger.Errorf("Refresh token generation failed for user ID %s: %v", user.ID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate refresh token"})
		return
	}

	// Hash the refresh token for storage
	hashedRefreshToken, err := util.HashRefreshToken(refreshTokenString)
	if err != nil {
		config.Logger.Errorf("Refresh token hashing failed for user ID %s: %v", user.ID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash refresh token"})
		return
	}

	// Create refresh token record in database
	expiresAt := time.Now().Add(30 * 24 * time.Hour) // 30 days
	refreshTokenRecord := models.RefreshToken{
		UserID:    user.ID,
		TokenHash: hashedRefreshToken,
		ExpiresAt: expiresAt,
		Revoked:   false,
	}

	if err := config.GetDB().Create(&refreshTokenRecord).Error; err != nil {
		config.Logger.Errorf("Failed to create refresh token record for user ID %s: %v", user.ID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create refresh token"})
		return
	}

	config.Logger.Infof("User login successful: ID %s, Email: %s", user.ID, user.Email)

	// Remove password from response
	user.Password = ""

	c.JSON(http.StatusOK, LoginResponse{
		Message:      "Login successful",
		AccessToken:  accessToken,
		RefreshToken: refreshTokenString,
		TokenType:    "Bearer",
		ExpiresIn:    900, // 15 minutes in seconds
		User:         user,
	})
}

// RegisterRequest represents the request body for user registration
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email" example:"newuser@example.com"`
	Name     string `json:"name" binding:"required" example:"John Doe"`
	Password string `json:"password" binding:"required" example:"securepassword123"`
}

// RegisterResponse represents the response body for successful registration
type RegisterResponse struct {
	Message      string    `json:"message" example:"Registration successful"`
	AccessToken  string    `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	RefreshToken string    `json:"refresh_token" example:"abc123def456..."`
	TokenType    string    `json:"token_type" example:"Bearer"`
	ExpiresIn    int       `json:"expires_in" example:"900"`
	UserID       uuid.UUID `json:"user_id" example:"550e8400-e29b-41d4-a716-446655440000"`
}

// Register godoc
// @Summary      User registration
// @Description  Register a new user and return JWT token
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user  body      RegisterRequest  true  "Registration data"
// @Success      201   {object}  RegisterResponse
// @Failure      400   {object}  map[string]string
// @Failure      409   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /auth/register [post]
func Register(c *gin.Context) {
	var input RegisterRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		config.Logger.Warnf("Invalid registration input: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	config.Logger.Infof("Registration attempt for email: %s, name: %s", input.Email, input.Name)

	// Check if user already exists
	var existingUser models.User
	if err := config.GetDB().Where("email = ?", input.Email).First(&existingUser).Error; err == nil {
		config.Logger.Warnf("Registration failed - email already exists: %s", input.Email)
		c.JSON(http.StatusConflict, gin.H{"error": "Email already registered"})
		return
	}

	hashedPassword, err := util.HashPassword(input.Password)
	if err != nil {
		config.Logger.Errorf("Password hashing failed for email %s: %v", input.Email, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process registration"})
		return
	}

	user := models.User{
		Email:    input.Email,
		Name:     input.Name,
		Password: hashedPassword,
		Settings: util.GetDefaultUserSettings(),
	}

	if err := config.GetDB().Create(&user).Error; err != nil {
		config.Logger.Errorf("Error creating user for email %s: %v", input.Email, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user account"})
		return
	}

	// Generate access token (short-lived)
	accessToken, err := util.GenerateAccessToken(user.ID)
	if err != nil {
		config.Logger.Errorf("Access token generation failed for new user ID %s: %v", user.ID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Registration completed but failed to generate access token"})
		return
	}

	// Generate refresh token (long-lived)
	refreshTokenString, err := util.GenerateRefreshToken()
	if err != nil {
		config.Logger.Errorf("Refresh token generation failed for new user ID %s: %v", user.ID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Registration completed but failed to generate refresh token"})
		return
	}

	// Hash the refresh token for storage
	hashedRefreshToken, err := util.HashRefreshToken(refreshTokenString)
	if err != nil {
		config.Logger.Errorf("Refresh token hashing failed for new user ID %s: %v", user.ID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Registration completed but failed to hash refresh token"})
		return
	}

	// Create refresh token record in database
	expiresAt := time.Now().Add(30 * 24 * time.Hour) // 30 days
	refreshTokenRecord := models.RefreshToken{
		UserID:    user.ID,
		TokenHash: hashedRefreshToken,
		ExpiresAt: expiresAt,
		Revoked:   false,
	}

	if err := config.GetDB().Create(&refreshTokenRecord).Error; err != nil {
		config.Logger.Errorf("Failed to create refresh token record for new user ID %s: %v", user.ID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Registration completed but failed to create refresh token"})
		return
	}

	config.Logger.Infof("User registered successfully: ID %s, Email: %s", user.ID, user.Email)
	c.JSON(http.StatusCreated, RegisterResponse{
		Message:      "Registration successful",
		AccessToken:  accessToken,
		RefreshToken: refreshTokenString,
		TokenType:    "Bearer",
		ExpiresIn:    900, // 15 minutes in seconds
		UserID:       user.ID,
	})
}

// UpdateUserRequest represents the request body for updating a user
type UpdateUserRequest struct {
	Name     *string                `json:"name" example:"Updated Name"`
	Email    *string                `json:"email" example:"updated@example.com"`
	Password *string                `json:"password" example:"newpassword123"`
	Settings map[string]interface{} `json:"settings" example:"{\"theme\": \"dark\"}"`
}

// UpdateUser godoc
// @Summary      Update user information
// @Description  Update user profile information
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        ID    path      int                true  "User ID"
// @Param        user  body      UpdateUserRequest  true  "User update data"
// @Success      200   {object}  models.User
// @Failure      400   {object}  map[string]string
// @Failure      401   {object}  map[string]string
// @Failure      403   {object}  map[string]string
// @Failure      404   {object}  map[string]string
// @Failure      409   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /users/{ID} [put]
func UpdateUser(c *gin.Context) {
	userIDStr := c.Param("ID")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		config.Logger.Warnf("Invalid user ID param for update: %s", userIDStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Check if the requesting user is updating their own profile or is admin
	requestingUserID, exist := c.Get("userID")
	if !exist {
		config.Logger.Warn("userID not found in context during user update")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	requestingUserIDUUID, ok := requestingUserID.(uuid.UUID)
	if !ok {
		config.Logger.Errorf("Invalid userID type in context: %T", requestingUserID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// Users can only update their own profile (unless they're admin)
	if requestingUserIDUUID != userID {
		config.Logger.Warnf("User %s attempted to update user %s (forbidden)", requestingUserIDUUID, userID)
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only update your own profile"})
		return
	}

	var user models.User
	if err := config.GetDB().First(&user, userID).Error; err != nil {
		config.Logger.Warnf("User not found for update: ID %s", userID)
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var input UpdateUserRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		config.Logger.Warnf("Invalid update input for user ID %d: %v", userID, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	updates := map[string]interface{}{}

	if input.Email != nil {
		// Check if email is already taken by another user
		var existingUser models.User
		if err := config.GetDB().Where("email = ? AND id != ?", *input.Email, userID).First(&existingUser).Error; err == nil {
			config.Logger.Warnf("Email update failed - email already exists: %s", *input.Email)
			c.JSON(http.StatusConflict, gin.H{"error": "Email already in use"})
			return
		}
		updates["email"] = *input.Email
	}

	if input.Name != nil {
		updates["name"] = *input.Name
	}

	if input.Password != nil {
		hashedPassword, err := util.HashPassword(*input.Password)
		if err != nil {
			config.Logger.Errorf("Password hashing failed for user ID %d: %v", userID, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
			return
		}
		updates["password"] = hashedPassword
	}

	if input.Settings != nil {
		updates["settings"] = input.Settings
	} else {
		config.Logger.Warnf("No valid fields provided for user update: ID %d", userID)
		c.JSON(http.StatusBadRequest, gin.H{"error": "No valid fields to update"})
		return
	}

	config.Logger.Infof("Updating user ID %d with fields: %v", userID, getUpdateFieldNames(updates))
	if err := config.GetDB().Model(&user).Updates(updates).Error; err != nil {
		config.Logger.Errorf("Error updating user ID %d: %v", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	// Reload user to get updated data and remove password from response
	if err := config.GetDB().First(&user, userID).Error; err != nil {
		config.Logger.Errorf("Error retrieving updated user ID %d: %v", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not reload updated user"})
		return
	}

	user.Password = "" // Remove password from response
	config.Logger.Infof("User updated successfully: ID %d", userID)
	c.JSON(http.StatusOK, user)
}

// DeleteUser godoc
// @Summary      Delete user account
// @Description  Delete a user account (self or admin)
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        ID   path      int  true  "User ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /users/{ID} [delete]
func DeleteUser(c *gin.Context) {
	userIDStr := c.Param("ID")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		config.Logger.Warnf("Invalid user ID param for delete: %s", userIDStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Check if the requesting user is deleting their own account or is admin
	requestingUserID, exist := c.Get("userID")
	if !exist {
		config.Logger.Warn("userID not found in context during user deletion")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	requestingUserIDUUID, ok := requestingUserID.(uuid.UUID)
	if !ok {
		config.Logger.Errorf("Invalid userID type in context: %T", requestingUserID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// Users can only delete their own account (unless they're admin)
	if requestingUserIDUUID != userID {
		config.Logger.Warnf("User %s attempted to delete user %s (forbidden)", requestingUserIDUUID, userID)
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only delete your own account"})
		return
	}

	var user models.User
	if err := config.GetDB().First(&user, userID).Error; err != nil {
		config.Logger.Warnf("User not found for deletion: ID %s", userID)
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	config.Logger.Infof("Deleting user account: ID %s, Email: %s", user.ID, user.Email)

	if err := config.GetDB().Delete(&user).Error; err != nil {
		config.Logger.Errorf("Failed to delete user ID %d: %v", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user account"})
		return
	}

	config.Logger.Infof("User account deleted successfully: ID %d", userID)
	c.JSON(http.StatusOK, gin.H{
		"message": "User account deleted successfully",
		"user_id": userID,
	})
}

// Helper function to get field names from updates map (for logging)
func getUpdateFieldNames(updates map[string]interface{}) []string {
	fields := make([]string, 0, len(updates))
	for field := range updates {
		if field == "password" {
			fields = append(fields, "password (hashed)")
		} else {
			fields = append(fields, field)
		}
	}
	return fields
}
