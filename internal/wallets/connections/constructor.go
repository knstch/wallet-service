package connections

import (
	"context"
	"fmt"
	"math/big"

	"github.com/knstch/subtrack-libs/log"

	"wallets-service/config"
	"wallets-service/internal/domain/dto"
	"wallets-service/internal/domain/enum"
	"wallets-service/internal/wallets/connections/blockchain"
)

type Blockchain interface {
	GetNativeBalance(ctx context.Context, walletAddr string, network enum.Network) (*big.Float, error)
	GetTokenBalanceAndInfo(ctx context.Context, walletAddr, tokenAddr string, network enum.Network) (dto.TokenInfo, error)
}

type Connections struct {
	Blockchain Blockchain
}

func MakeConnections(cfg *config.Config, lg *log.Logger) (*Connections, error) {
	blockchainClient, err := blockchain.NewClient(cfg, lg)
	if err != nil {
		return nil, fmt.Errorf("blockchain.NewClient: %w", err)
	}

	return &Connections{
		Blockchain: blockchainClient,
	}, nil
}
