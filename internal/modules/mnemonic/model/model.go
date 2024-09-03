package model

import "nexus-wallet/pkg/repository"

type Mnemonic struct {
	Id   int64  `primary:"true" must_generate:"true" db:"id"`
	Name string `db:"name"`
}

func (Mnemonic) GetTableName() string {
	return "mnemonics"
}

func (Mnemonic) Clear() repository.Model {
	return &Mnemonic{}
}
