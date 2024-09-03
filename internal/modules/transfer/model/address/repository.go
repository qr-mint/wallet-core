package address

import (
	"database/sql/driver"
	"github.com/doug-martin/goqu/v9"
	"nexus-wallet/internal/app_enum"
	"nexus-wallet/pkg/repository"
)

type Repository struct {
	*repository.BaseRepository
}

func NewRepository(baseRepository *repository.BaseRepository) *Repository {
	return &Repository{baseRepository}
}

type FindOptions struct {
	MnemonicId int64
	Network    app_enum.Network
}

func (r Repository) Find(options FindOptions, tx driver.Tx) (*Address, error) {
	address := Address{}
	err := r.FindOneBy(goqu.Ex{"mnemonic_id": options.MnemonicId, "network": options.Network}, &address, tx)
	if err != nil {
		return nil, err
	}

	return &address, nil
}
