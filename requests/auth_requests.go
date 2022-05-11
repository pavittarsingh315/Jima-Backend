package requests

type RegistrationRequest struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
	Contact  string `json:"contact"`
	Code     int    `json:"code"`
}

type LoginRequest struct {
	Contact  string `json:"contact"`
	Password string `json:"password"`
}

type TokenLoginRequest struct {
	AccessToken  string `json:"access"`
	RefreshToken string `json:"refresh"`
}
