package connections

import (
	"fmt"
	"github.com/knstch/subtrack-libs/log"
	"wallets-service/config"
	"wallets-service/internal/wallets/connections/blockchain-gateway"
)

type Connections struct {
	Blockchain blockchain.Blockchain
}

func MakeConnections(cfg config.Config, lg *log.Logger) (*Connections, error) {
	blockchainGateway, err := blockchain.MakeBlockchainGatewayClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("blockchain.MakeBlockchainGatewayClient: %w", err)
	}

	blockchainGatewayClient := blockchain.NewClient(lg, blockchainGateway)

	return &Connections{
		Blockchain: blockchainGatewayClient,
	}, nil
}
