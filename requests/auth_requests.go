package requests

type RegistrationRequest struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
	Contact  string `json:"contact"`
	Code     string `json:"code"`
}

type LoginRequest struct {
	Contact  string `json:"contact"`
	Password string `json:"password"`
}

type TokenLoginRequest struct {
	AccessToken  string `json:"access"`
	RefreshToken string `json:"refresh"`
}

type RecoveryRequest struct {
	Contact  string `json:"contact"`
	Password string `json:"password"`
	Code     string `json:"code"`
}
