package request

// LoginRequest
// @Description LoginRequest is a struct that represents the request of login
type LoginRequest struct {
	UserID   string `json:"user_id"`
	Password string `json:"password"`
}
