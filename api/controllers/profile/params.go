package profile

import "nexus-wallet/internal/modules/profile"

type GetProfileTelegramResponse struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Username    string `json:"username"`
	Language    string `json:"language"`
	ImageSource string `json:"image_source"`
}

func (GetProfileTelegramResponse) fillFromOutput(output profile.GetOutput) GetProfileTelegramResponse {
	return GetProfileTelegramResponse{
		FirstName:   output.TelegramProfile.FirstName,
		LastName:    output.TelegramProfile.LastName,
		Username:    output.TelegramProfile.Username,
		ImageSource: output.TelegramProfile.ImageSource,
		Language:    string(output.Language),
	}
}
