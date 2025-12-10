package routes

// LoginRequest represents the login credentials
// @Description Login credentials for authentication
type LoginRequest struct {
	User     string `json:"user" binding:"required" example:"admin"`
	Password string `json:"password" binding:"required" example:"password123"`
}

// UpdatePasswordRequest represents the password update request
// @Description Request to update user password
type UpdatePasswordRequest struct {
	CurrentPassword string `json:"currentPassword" binding:"required" example:"oldPassword123"`
	NewPassword     string `json:"newPassword" binding:"required" validate:"gte=8,lte=64" example:"newPassword123"`
}

// ErrorResponse represents an error response
// @Description Error response
type ErrorResponse struct {
	Error string `json:"error" example:"error message"`
}

// SuccessResponse represents a success response
// @Description Success response
type SuccessResponse struct {
	Message string `json:"message" example:"operation successful"`
}

// CreateUserRequest represents the request to create a new user
// @Description Request to create a new user (admin only)
type CreateUserRequest struct {
	Name     string `json:"name" binding:"required" validate:"required,gte=3,lte=32" example:"john"`
	Password string `json:"password" binding:"required" validate:"required,gte=8,lte=64" example:"password123"`
	Admin    bool   `json:"admin" example:"false"`
}

// UpdateUserRequest represents the request to update a user
// @Description Request to update a user (admin only)
type UpdateUserRequest struct {
	Admin    *bool   `json:"admin,omitempty" example:"false"`
	Password *string `json:"password,omitempty" validate:"omitempty,gte=8,lte=64" example:"newPassword123"`
}
