package wallets

import (
	"context"
	"fmt"

	"github.com/knstch/subtrack-libs/auth"

	"wallets-service/internal/domain/dto"
	"wallets-service/internal/domain/enum"
)

func (svc *ServiceImpl) GetBalance(ctx context.Context, network enum.Network) {
	userInfo, err := auth.GetUserData(ctx)
	if err != nil {
		return nil, fmt.Errorf("auth.GetUserInfoFromContext: %w", err)
	}

	wallet, err := svc.repo.GetWalletByUserID(ctx, userInfo.UserID)
	if err != nil {
		return nil, fmt.Errorf("repo.GetWalletByUserID: %w", err)
	}

	nativeBalance, err := svc.conns.Blockchain.GetNativeBalance(ctx, wallet.PublicAddr, network)
	if err != nil {
		return nil, fmt.Errorf("connections.GetNativeBalance: %w", err)
	}

	tokens := make([]dto.TokenInfo, 0, len(wallet.TokenAddresses))
	for _, address := range wallet.TokenAddresses {
		token, err := svc.conns.Blockchain.GetTokenBalanceAndInfo(ctx, wallet.PublicAddr, address, network)
		if err != nil {
			return nil, fmt.Errorf("connections.GetTokenBalanceAndInfo: %w", err)
		}

		tokens = append(tokens, token)
	}
}
