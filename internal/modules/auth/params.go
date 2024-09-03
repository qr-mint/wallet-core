package auth

type CheckAccessTokenInput struct {
	AccessToken string
}

type CheckAccessTokenOutput struct {
	UserId int64
}

type RefreshInput struct {
	RefreshToken string
}

type AuthenticateOutput struct {
	AccessToken  string
	RefreshToken string
}
