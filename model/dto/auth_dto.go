package dto

type AuthRequestDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponseDTO struct {
	Token string `json:"token"`
}
