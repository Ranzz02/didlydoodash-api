package dto

// ---- Structs ----
type UserResponse struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

type Tokens struct {
	UserID  string `json:"-"`
	Access  string `json:"access"`
	Refresh string `json:"refresh"`
}

// ---- Request Structs ----
type SignInRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	Remember bool   `json:"remember" default:"false"`
}

type SignUpRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	Remember bool   `json:"remember" default:"false"`
}

type RefreshRequest struct {
	Token string `json:"token" binding:"required"`
}

// ---- Response Structs ----
type AuthResponse struct {
	User   UserResponse `json:"user"`
	Tokens Tokens       `json:"tokens"`
}
