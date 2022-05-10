package requests

type RegistrationRequest struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
	Contact  string `json:"contact"`
	Code     int    `json:"code"`
}
