package blockchain

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/knstch/subtrack-libs/log"

	"wallets-service/config"
	"wallets-service/internal/domain/dto"
	"wallets-service/internal/domain/enum"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/ethclient"
)

const erc20ABIJSON = `[
  {"constant":true,"inputs":[{"name":"account","type":"address"}],"name":"balanceOf","outputs":[{"name":"","type":"uint256"}],"type":"function"},
  {"constant":true,"inputs":[],"name":"name","outputs":[{"name":"","type":"string"}],"type":"function"},
  {"constant":true,"inputs":[],"name":"symbol","outputs":[{"name":"","type":"string"}],"type":"function"},
  {"constant":true,"inputs":[],"name":"decimals","outputs":[{"name":"","type":"uint8"}],"type":"function"}
]`

type ClientImpl struct {
	PolygonAddr string
	erc20ABI    abi.ABI
	lg          *log.Logger
}

func NewClient(cfg *config.Config, logger *log.Logger) (*ClientImpl, error) {
	erc20ABI, err := abi.JSON(strings.NewReader(erc20ABIJSON))
	if err != nil {
		return nil, err
	}

	return &ClientImpl{
		PolygonAddr: cfg.Blockchains.PolygonAddr,
		erc20ABI:    erc20ABI,
		lg:          logger,
	}, nil
}

func (c *ClientImpl) getRpcUrl(network enum.Network) string {
	switch network {
	case enum.PolygonNetwork:
		return c.PolygonAddr
	default:
		return ""
	}
}

func (c *ClientImpl) GetNativeBalance(ctx context.Context, walletAddr string, network enum.Network) (*big.Float, error) {
	client, err := ethclient.Dial(c.getRpcUrl(network))
	if err != nil {
		return nil, fmt.Errorf("ethclient.Dial: %w", err)
	}
	defer client.Close()

	account := common.HexToAddress(walletAddr)

	balanceWei, err := client.BalanceAt(ctx, account, nil)
	if err != nil {
		return nil, fmt.Errorf("client.BalanceAt: %w", err)
	}

	balanceCurrency := new(big.Float)
	balanceCurrency.SetString(balanceWei.String())
	ethValue := new(big.Float).Quo(balanceCurrency, big.NewFloat(1e18))

	return ethValue, nil
}

func (c *ClientImpl) GetTokenBalanceAndInfo(ctx context.Context, walletAddr, tokenAddr string, network enum.Network) (dto.TokenInfo, error) {
	client, err := ethclient.Dial(c.getRpcUrl(network))
	if err != nil {
		return dto.TokenInfo{}, fmt.Errorf("ethclient.Dial: %w", err)
	}
	defer client.Close()

	walletAddressCommon := common.HexToAddress(walletAddr)
	tokenAddressCommon := common.HexToAddress(tokenAddr)

	contract := bind.NewBoundContract(tokenAddressCommon, c.erc20ABI, client, client, client)

	balance, err := getBalance(ctx, walletAddressCommon, contract)
	if err != nil {
		return dto.TokenInfo{}, err
	}

	symbol, err := getSymbol(ctx, contract)
	if err != nil {
		return dto.TokenInfo{}, err
	}

	return dto.TokenInfo{
		Balance: balance,
		Symbol:  symbol,
	}, nil
}

func getBalance(ctx context.Context, walletAddr common.Address, contract *bind.BoundContract) (*big.Float, error) {
	var outBalance []interface{}
	err := contract.Call(&bind.CallOpts{Context: ctx}, &outBalance, "balanceOf", walletAddr)
	if err != nil {
		return nil, fmt.Errorf("contract.Call: %w", err)
	}

	balance, ok := outBalance[0].(*big.Int)
	if !ok {
		return nil, fmt.Errorf("unexpected type in balanceOf result: %T", outBalance[0])
	}

	readableBalance := new(big.Float).Quo(new(big.Float).SetInt(balance), big.NewFloat(1e18))

	return readableBalance, nil
}

func getSymbol(ctx context.Context, contract *bind.BoundContract) (string, error) {
	var outSymbol []interface{}
	err := contract.Call(&bind.CallOpts{Context: ctx}, &outSymbol, "symbol")
	if err != nil {
		return "", fmt.Errorf("contract.Call: %w", err)
	}
	symbol, ok := outSymbol[0].(string)
	if !ok {
		symbol = ""
	}

	return symbol, nil
}
