package mnemonic

type Mnemonic struct {
	Id   int64  `primary:"true" must_generate:"true" db:"id"`
	Name string `db:"name"`
	Hash string `db:"hash"`
}

func (Mnemonic) GetTableName() string {
	return "mnemonics"
}

type UsersMnemonics struct {
	Id         int64 `primary:"true" must_generate:"true" db:"id"`
	UserId     int64 `db:"user_id"`
	MnemonicId int64 `db:"mnemonic_id"`
}

func (UsersMnemonics) GetTableName() string {
	return "users_mnemonics"
}
