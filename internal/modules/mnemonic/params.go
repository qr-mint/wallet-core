package mnemonic

type GetMnemonicIdInput struct {
	Hash string
}

type GetMnemonicIdOutput struct {
	Id int64
}

type GenerateInput struct {
	UserId int64
}

type GenerateOutput struct {
	Mnemonic string
}

type UpdateNameInput struct {
	Name       string
	MnemonicId int64
}

type GetNamesInput struct {
	UserId int64
}

type GetNamesOutputItem struct {
	Name string
}

type GetNamesOutput struct {
	Items []GetNamesOutputItem
}
