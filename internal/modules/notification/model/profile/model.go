package profile

import (
	"nexus-wallet/internal/app_enum"
	"nexus-wallet/pkg/repository"
)

type Profile struct {
	Id       int64                `primary:"true" must_generate:"true" db:"id"`
	Language app_enum.Language    `db:"language"`
	Type     app_enum.ProfileType `db:"type"`
	UserId   int64                `db:"user_id"`
}

func (p *Profile) GetTableName() string {
	return "profiles"
}

func (p *Profile) Clear() repository.Model {
	return &Profile{}
}
