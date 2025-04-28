package wallets

type WalletWithBalance struct {
	NativeBalance string  `json:"native_balance"`
	Tokens        []Token `json:"tokens"`
}

type Token struct {
	Balance string `json:"balance"`
	Symbol  string `json:"symbol"`
}

type Wallet struct {
	PublicKey  string
	PrivateKey []byte
}
