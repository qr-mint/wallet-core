package user

import (
	"database/sql/driver"
	"fmt"
	tgbotAPI "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
	"gitlab.com/golib4/logger/logger"
	"nexus-wallet/internal/app_enum"
	"nexus-wallet/internal/app_error"
	"nexus-wallet/internal/modules/auth/model/profile"
	user_model "nexus-wallet/internal/modules/auth/model/user"
	"nexus-wallet/pkg/repository"
	"nexus-wallet/pkg/transaction"
)

type TelegramUserFactory struct {
	telegramClient            *tgbotAPI.BotAPI
	profileRepository         *profile.Repository
	telegramProfileRepository *profile.TelegramProfileRepository
	telegramUserRepository    *user_model.TelegramUserRepository
	logger                    logger.Logger
	transactionManager        transaction.Manager
	botToken                  string
}

func NewTelegramUserFactory(
	telegramClient *tgbotAPI.BotAPI,
	profileRepository *profile.Repository,
	telegramProfileRepository *profile.TelegramProfileRepository,
	telegramUserRepository *user_model.TelegramUserRepository,
	transactionManager transaction.Manager,
	logger logger.Logger,
	botToken string,
) *TelegramUserFactory {
	return &TelegramUserFactory{
		telegramClient:            telegramClient,
		profileRepository:         profileRepository,
		telegramProfileRepository: telegramProfileRepository,
		telegramUserRepository:    telegramUserRepository,
		transactionManager:        transactionManager,
		logger:                    logger,
		botToken:                  botToken,
	}
}

type CreateData struct {
	TelegramId   int64
	FirstName    string
	LastName     string
	Username     string
	LanguageCode string
}

func (f TelegramUserFactory) Create(data CreateData) (int64, *app_error.AppError) {
	var telegramUser *user_model.TelegramUser
	var telegramProfile *profile.TelegramProfile
	var transactionErr error
	var findUserErr error
	err := f.transactionManager.WithTransaction(func(tx driver.Tx) error {
		telegramUser, transactionErr = f.createUser(data, tx)
		if transactionErr != nil {
			if repository.IsUniqueError(transactionErr) {
				telegramUser, telegramProfile, findUserErr = f.findExistingUser(data, nil)
				if findUserErr != nil {
					return fmt.Errorf("can not get existing telegram user: %s", transactionErr)
				}
				if telegramUser == nil {
					return errors.New("unique exception was catch but user still not found on user creation")
				}

				return transactionErr
			}

			return fmt.Errorf("can not create telegram user: %s", transactionErr)
		}
		if telegramProfile == nil {
			telegramProfile, transactionErr = f.createProfile(data, *telegramUser, tx)
			if transactionErr != nil {
				return fmt.Errorf("can not create telegram profile: %s", transactionErr)
			}
		}

		return nil
	})
	if err != nil && !repository.IsUniqueError(err) {
		return 0, app_error.InternalError(fmt.Errorf("error in transaction: %s", err))
	}

	go func() {
		err := f.resolveImageSource(telegramProfile, telegramUser.TelegramId)
		if err != nil {
			f.logger.Warningf("can not resolve telegram image source: %s", err)
		}
	}()

	return telegramUser.Id, nil
}

func (f TelegramUserFactory) findExistingUser(data CreateData, tx driver.Tx) (*user_model.TelegramUser, *profile.TelegramProfile, error) {
	telegramUser, err := f.telegramUserRepository.Find(user_model.FindOptions{TelegramID: data.TelegramId}, tx)
	if err != nil {
		return nil, nil, fmt.Errorf("can not get telegram user: %s", err)
	}
	if telegramUser != nil {
		telegramProfile, err := f.telegramProfileRepository.Find(profile.FindOptions{UserId: telegramUser.Id}, tx)
		if err != nil {
			return nil, nil, fmt.Errorf("can not get telegram profile: %s", err)
		}
		if telegramProfile == nil {
			return nil, nil, fmt.Errorf("can not find telegram profile by userId: %d", telegramUser.Id)
		}

		return telegramUser, telegramProfile, nil
	}

	return nil, nil, nil
}

func (f TelegramUserFactory) createUser(data CreateData, tx driver.Tx) (*user_model.TelegramUser, error) {
	baseUser := user_model.User{Type: app_enum.TelegramUserType}
	userModel := &user_model.TelegramUser{TelegramId: data.TelegramId, TelegramBotToken: f.botToken, User: baseUser}
	err := f.telegramUserRepository.Save(userModel, tx)
	if err != nil {
		return nil, fmt.Errorf("can not save telegram user: %s", err)
	}

	return userModel, nil
}

func (f TelegramUserFactory) createProfile(data CreateData, userModel user_model.TelegramUser, tx driver.Tx) (*profile.TelegramProfile, error) {
	languageCode := app_enum.ToLanguage(data.LanguageCode)
	if languageCode == nil {
		engLanguageCode := app_enum.EngLanguage
		languageCode = &engLanguageCode
	}

	baseProfile := profile.Profile{Language: *languageCode, Type: app_enum.TelegramProfileType, UserId: userModel.UserId}
	profileModel := &profile.TelegramProfile{
		FirstName:   data.FirstName,
		LastName:    data.LastName,
		Username:    data.Username,
		ImageSource: "",
		Profile:     baseProfile,
	}
	err := f.profileRepository.Save(profileModel, tx)
	if err != nil {
		return nil, fmt.Errorf("can not create telegram profile: %s", err)
	}

	return profileModel, nil
}

func (f TelegramUserFactory) resolveImageSource(profileModel *profile.TelegramProfile, telegramId int64) error {
	photos, err := f.telegramClient.GetUserProfilePhotos(tgbotAPI.NewUserProfilePhotos(telegramId))
	if err != nil {
		return fmt.Errorf("can not get telegram user photo: %s", err)
	}

	fileURL := ""
	if photos.TotalCount > 0 {
		fileID := photos.Photos[0][len(photos.Photos[0])-1].FileID
		fileURL, err = f.telegramClient.GetFileDirectURL(fileID)
		if err != nil {
			return fmt.Errorf("can not get telegram user photo direct url: %s", err)
		}
	}

	profileModel.ImageSource = fileURL
	err = f.profileRepository.Save(profileModel, nil)
	if err != nil {
		return fmt.Errorf("can not save telegram profile image url: %s", err)
	}

	return nil
}
