package blockchain

import (
	"fmt"
	"wallets-service/config"

	blockchainGatewayApi "github.com/knstch/blockchain-gateway-api/private"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func MakeBlockchainGatewayClient(cfg config.Config) (blockchainGatewayApi.BlockchainGatewayClient, error) {
	conn, err := grpc.NewClient(cfg.BlockchainGatewayHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("grpc.NewClient: %w", err)
	}

	return blockchainGatewayApi.NewBlockchainGatewayClient(conn), nil
}
