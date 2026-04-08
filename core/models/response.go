package models

type LoginResponse struct {
	Success bool     `json:"success"`
	Token   string   `json:"token,omitempty"`
	User    UserInfo `json:"user,omitempty"`
	Error   string   `json:"error,omitempty"`
}

type UserInfo struct {
	ID        int32  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
}
