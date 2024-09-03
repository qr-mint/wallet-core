package transaction

import (
	"fmt"
	"nexus-wallet/internal/app_error"
	"nexus-wallet/internal/modules/transaction/model/coin"
	"nexus-wallet/internal/modules/transaction/model/transaction"
	"nexus-wallet/internal/modules/transaction/provider"
)

type Service struct {
	transactionRepository *transaction.Repository
	coinRepository        *coin.Repository
	syncer                *provider.Syncer
}

func NewService(
	transactionRepository *transaction.Repository,
	coinRepository *coin.Repository,
	syncer *provider.Syncer,
) *Service {
	return &Service{
		transactionRepository: transactionRepository,
		coinRepository:        coinRepository,
		syncer:                syncer,
	}
}

func (s Service) List(input ListInput) (*ListOutput, *app_error.AppError) {
	err := s.syncer.Sync(input.MnemonicId)
	if err != nil {
		return nil, app_error.InternalError(fmt.Errorf("can not sync transactions: %s", err))
	}

	transactionList, err := s.transactionRepository.FindMany(
		transaction.FindManyOptions{
			AddressCoinId: input.AddressCoinId,
			OnlyOut:       input.OnlyOut,
			MnemonicId:    input.MnemonicId,
			Limit:         input.Limit,
			Offset:        input.Offset,
		},
		nil,
	)
	if err != nil {
		return nil, app_error.InternalError(fmt.Errorf("can not get transactions list: %s", err))
	}

	coinsMappedByIds, err := s.coinRepository.FindAllMappedByIds(nil)
	if err != nil {
		return nil, app_error.InternalError(fmt.Errorf("can not get mapped coins by ids: %s", err))
	}

	output, err := ListOutput{}.fillFromModel(transactionList, coinsMappedByIds)
	if err != nil {
		return nil, app_error.InternalError(fmt.Errorf("can not fill transactions output: %s", err))
	}

	return output, nil
}
