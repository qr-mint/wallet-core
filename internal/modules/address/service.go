package address

import (
	"fmt"
	"github.com/pkg/errors"
	"nexus-wallet/internal/app_error"
	"nexus-wallet/internal/app_util"
	"nexus-wallet/internal/modules/address/balance"
	"nexus-wallet/internal/modules/address/import"
	"nexus-wallet/internal/modules/address/model/address"
	"nexus-wallet/internal/modules/address/model/coin"
	"nexus-wallet/internal/modules/address/model/mnemonic"
)

type Service struct {
	addressCoinRepository *address.AddressCoinRepository
	addressRepository     *address.AddressRepository
	importer              *_import.Importer
	coinRepository        *coin.CoinRepository
	priceRepository       *coin.PriceRepository
	mnemonicRepository    *mnemonic.Repository
	syncer                *balance.Syncer
}

func NewService(
	addressCoinRepository *address.AddressCoinRepository,
	addressRepository *address.AddressRepository,
	importer *_import.Importer,
	coinRepository *coin.CoinRepository,
	priceRepository *coin.PriceRepository,
	mnemonicRepository *mnemonic.Repository,
	syncer *balance.Syncer,
) *Service {
	return &Service{
		addressCoinRepository: addressCoinRepository,
		addressRepository:     addressRepository,
		importer:              importer,
		coinRepository:        coinRepository,
		priceRepository:       priceRepository,
		mnemonicRepository:    mnemonicRepository,
		syncer:                syncer,
	}
}

func (s Service) Import(input ImportInput) *app_error.AppError {
	return s.importer.Import(input.UserId, input.ImportData)
}

func (s Service) SwitchCoinVisibility(input SwitchCoinVisibilityInput) *app_error.AppError {
	findOptions := address.FindOptions{Id: input.AddressCoinId, MnemonicId: input.MnemonicId}
	addressCoin, err := s.addressCoinRepository.Find(findOptions, nil)
	if err != nil {
		return app_error.InternalError(err)
	}

	addressCoin.IsVisible = !addressCoin.IsVisible
	err = s.addressCoinRepository.Save(addressCoin, nil)
	if err != nil {
		return app_error.InternalError(err)
	}

	return nil
}

func (s Service) GetCoin(input GetCoinInput) (*GetCoinOutput, *app_error.AppError) {
	addressCoin, err := s.addressCoinRepository.Find(address.FindOptions{Id: input.Id, MnemonicId: input.MnemonicId}, nil)
	if err != nil {
		return nil, app_error.InternalError(fmt.Errorf("find address coin error in address service: %s", err))
	}
	if addressCoin == nil {
		return nil, app_error.InvalidDataError(errors.Errorf("address coin with id %d not found", input.Id))
	}
	coinsMappedByIds, err := s.coinRepository.FindAllMappedByIds(nil)
	if err != nil {
		return nil, app_error.InternalError(fmt.Errorf("can not get mapped coins by ids in wallet service: %s", err))
	}
	addressesMappedByIds, err := s.addressRepository.FindMany(address.FindManyAddressOptions{MnemonicId: input.MnemonicId}, nil)
	if err != nil {
		return nil, app_error.InternalError(fmt.Errorf("can not get mapped addresses in address service: %s", err))
	}

	result, err := GetCoinOutput{}.fillFromModel(coinsMappedByIds[addressCoin.CoinId], *addressCoin, addressesMappedByIds[addressCoin.AddressId])
	if err != nil {
		return nil, app_error.InternalError(fmt.Errorf("can not fill coin list item from model"))
	}

	return result, nil
}

func (s Service) GetCoinsList(input GetCoinsListInput) (*GetCoinsListOutput, *app_error.AppError) {
	findOptions := address.FindManyOptions{MnemonicId: input.MnemonicId, OnlyVisible: input.OnlyVisible}
	addressCoins, err := s.addressCoinRepository.FindMany(findOptions, nil)
	if err != nil {
		return nil, app_error.InternalError(fmt.Errorf("find list address coins error in address service: %s", err))
	}

	coinsMappedByIds, err := s.coinRepository.FindAllMappedByIds(nil)
	if err != nil {
		return nil, app_error.InternalError(fmt.Errorf("can not get mapped coins by ids in wallet service: %s", err))
	}

	addressesMappedByIds, err := s.addressRepository.FindMany(address.FindManyAddressOptions{MnemonicId: input.MnemonicId}, nil)
	if err != nil {
		return nil, app_error.InternalError(fmt.Errorf("can not get mapped addresses in address service: %s", err))
	}

	var outputItems []GetCoinListOutputItem
	for _, addressCoin := range addressCoins {
		item, err := GetCoinListOutputItem{}.fillFromModel(
			coinsMappedByIds[addressCoin.CoinId],
			*addressCoin,
			addressesMappedByIds[addressCoin.AddressId],
		)
		if err != nil {
			return nil, app_error.InternalError(fmt.Errorf("can not fill coin list item from model"))
		}
		outputItems = append(outputItems, *item)
	}

	return &GetCoinsListOutput{Items: outputItems}, nil
}

func (s Service) GetAggregatedInfo(input GetAggregatedInfoInput) (*GetAffrefatedInfoOutput, *app_error.AppError) {
	err := s.syncer.Sync(input.MnemonicId)
	if err != nil {
		return nil, app_error.InternalError(errors.Errorf("can not sync balance: %s", err))
	}
	coinsList, appErr := s.GetCoinsList(GetCoinsListInput{MnemonicId: input.MnemonicId, OnlyVisible: true})
	if appErr != nil {
		return nil, appErr
	}
	var commonFiatAmount float64
	var commonDailyPriceDeltaPercent float64 = 0
	var outputItems []InfoOutputItem
	for _, coinData := range coinsList.Items {
		aggregatedCoinDataItem, err := s.getCoinAggregatedData(coinData, input)
		if err != nil {
			return nil, app_error.InternalError(errors.Errorf("can not get aggregatedCoinDataItem: %s", err))
		}

		commonDailyPriceDeltaPercent += aggregatedCoinDataItem.DailyPriceDeltaPercent
		commonFiatAmount += aggregatedCoinDataItem.FiatAmount
		outputItems = append(outputItems, *aggregatedCoinDataItem)
	}

	mnemonicData, err := s.mnemonicRepository.FindOne(input.MnemonicId, nil)
	if err != nil {
		return nil, app_error.InternalError(fmt.Errorf("can not find mnemonic: %s", err))
	}

	return &GetAffrefatedInfoOutput{
		DailyPriceDeltaPercent: commonDailyPriceDeltaPercent,
		FiatAmount:             commonFiatAmount,
		Currency:               input.Currency,
		Items:                  outputItems,
		Name:                   mnemonicData.Name,
	}, nil
}

func (s Service) getCoinAggregatedData(addressCoin GetCoinListOutputItem, data GetAggregatedInfoInput) (*InfoOutputItem, error) {
	priceData, err := s.priceRepository.FindLatest(coin.FindOptions{CoinId: addressCoin.CoinId, FiatCurrency: data.Currency}, nil)
	if err != nil {
		return nil, errors.Errorf("can not get price data: %s", err)
	}
	yesterdayPriceData, err := s.priceRepository.FindYesterday(coin.FindOptions{CoinId: addressCoin.CoinId, FiatCurrency: data.Currency}, nil)
	if err != nil {
		return nil, errors.Errorf("can not get yesterdayPriceData: %s", err)
	}
	var dailyPriceDeltaPercent float64 = 0
	if addressCoin.Amount != 0 && yesterdayPriceData != nil {
		dailyPriceDeltaPercent = ((float64(priceData.Price) - float64(yesterdayPriceData.Price)) / float64(yesterdayPriceData.Price)) * 100
	}

	var floatPrice float64
	if priceData != nil {
		floatPrice, err = app_util.AmountToFloat(addressCoin.Network, priceData.Price)
		if err != nil {
			return nil, fmt.Errorf("can not format price to float: %s", err)
		}
	}

	return &InfoOutputItem{
		Id:                     addressCoin.Id,
		FiatAmount:             addressCoin.Amount * floatPrice,
		Currency:               data.Currency,
		Amount:                 addressCoin.Amount,
		Symbol:                 addressCoin.Symbol,
		ImageSource:            addressCoin.ImageSource,
		Network:                addressCoin.Network,
		Name:                   addressCoin.Name,
		DailyPriceDeltaPercent: dailyPriceDeltaPercent,
	}, nil
}
