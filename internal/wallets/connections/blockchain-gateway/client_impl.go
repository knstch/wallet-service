package blockchain

import (
	"context"
	"fmt"
	blockchainGatewayApi "github.com/knstch/blockchain-gateway-api/private"
	"github.com/knstch/subtrack-libs/log"
	"github.com/knstch/subtrack-libs/tracing"

	"github.com/knstch/subtrack-libs/enum"
)

type Blockchain interface {
	GetBalance(ctx context.Context, pubAddr string, tokens []string, network enum.Network) (Balance, error)
}

type ClientImpl struct {
	conn blockchainGatewayApi.BlockchainGatewayClient
	lg   *log.Logger
}

func NewClient(lg *log.Logger, conn blockchainGatewayApi.BlockchainGatewayClient) *ClientImpl {
	return &ClientImpl{
		conn: conn,
		lg:   lg,
	}
}

func (c *ClientImpl) GetBalance(ctx context.Context, pubAddr string, tokens []string, network enum.Network) (Balance, error) {
	ctx, span := tracing.StartSpan(ctx, "client: GetBalance")
	defer span.End()

	resp, err := c.conn.GetBalance(ctx, &blockchainGatewayApi.GetBalanceRequest{
		PublicAddress:  pubAddr,
		TokenAddresses: tokens,
		Network:        network.String(),
	})
	if err != nil {
		return Balance{}, fmt.Errorf("conn.GetBalance: %w", err)
	}

	tokensBalance := make([]TokenBalance, 0, len(resp.Tokens))
	for _, token := range resp.Tokens {
		tokensBalance = append(tokensBalance, TokenBalance{
			Balance: token.Balance,
			Symbol:  token.Symbol,
		})
	}

	return Balance{
		NativeBalance: resp.NativeBalance,
		Tokens:        tokensBalance,
	}, nil
}
