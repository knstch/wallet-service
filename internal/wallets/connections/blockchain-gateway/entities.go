package blockchain

type Balance struct {
	NativeBalance string
	Tokens        []TokenBalance
}

type TokenBalance struct {
	Balance string
	Symbol  string
}
