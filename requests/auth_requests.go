package requests

type InitiateRegistrationRequest struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
	Contact  string `json:"contact"`
}
