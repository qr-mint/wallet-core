package profile

import (
	"errors"
	"github.com/gin-gonic/gin"
	"nexus-wallet/api/error_handler"
	"nexus-wallet/api/response"
	"nexus-wallet/internal/app_enum"
	"nexus-wallet/internal/app_error"
	"nexus-wallet/internal/modules/profile"
)

type ProfileController struct {
	profileService  *profile.Service
	responseFactory *response.ResponseFactory
	errorHandler    *error_handler.HttpErrorHandler
}

func NewProfileController(profileService *profile.Service, responseFactory *response.ResponseFactory) *ProfileController {
	return &ProfileController{
		profileService:  profileService,
		responseFactory: responseFactory,
	}
}

func (c ProfileController) GetProfile(context *gin.Context) {
	profileData, err := c.profileService.Get(profile.GetInput{UserId: context.GetInt64("userId")})
	if err != nil {
		c.errorHandler.Handle(context, err)
		return
	}
	if profileData.Type == app_enum.TelegramProfileType {
		c.responseFactory.Ok(context, GetProfileTelegramResponse{}.fillFromOutput(*profileData))
		return
	}

	c.errorHandler.Handle(context, app_error.InternalError(errors.New("unknown profile type returned")))
}
