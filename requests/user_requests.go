package requests

type EditProfileRequest struct {
	Username          string `json:"username"`
	Name              string `json:"name"`
	Bio               string `json:"bio"`
	BlacklistMessage  string `json:"blacklistMessage"`
	NewProfilePicture string `json:"newProfilePicture"`
	OldProfilePicture string `json:"oldProfilePicture"`
}
