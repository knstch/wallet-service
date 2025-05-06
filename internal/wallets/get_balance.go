package wallets

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/knstch/subtrack-libs/auth"
	"github.com/knstch/subtrack-libs/tracing"
	"time"

	"github.com/knstch/subtrack-libs/enum"
)

func (svc *ServiceImpl) GetBalance(ctx context.Context, network enum.Network) (*WalletWithBalance, error) {
	ctx, span := tracing.StartSpan(ctx, "service: GetBalance")
	defer span.End()

	userInfo, err := auth.GetUserData(ctx)
	if err != nil {
		return nil, fmt.Errorf("auth.GetUserInfoFromContext: %w", err)
	}

	balance := svc.getBalanceFromCache(userInfo.UserID)
	if balance != nil {
		return balance, nil
	}

	wallet, err := svc.repo.GetWalletInfoByUserID(ctx, userInfo.UserID)
	if err != nil {
		return nil, fmt.Errorf("repo.GetWalletInfoByUserID: %w", err)
	}

	balanceFromNetwork, err := svc.connections.Blockchain.GetBalance(ctx, wallet.PublicAddr, wallet.TokenAddresses, network)
	if err != nil {
		return nil, fmt.Errorf("svc.connections.Blockchain.GetBalance: %w", err)
	}

	balance = &WalletWithBalance{
		NativeBalance: balanceFromNetwork.NativeBalance,
	}

	balance.Tokens = make([]Token, 0, len(wallet.TokenAddresses))
	for _, v := range balanceFromNetwork.Tokens {
		balance.Tokens = append(balance.Tokens, Token{
			Balance: v.Balance,
			Symbol:  v.Symbol,
		})
	}

	svc.putBalanceToCache(userInfo.UserID, balance)

	return balance, nil
}

func getBalanceKey(userID uint) string {
	return fmt.Sprintf("balance:%d", userID)
}

func (svc *ServiceImpl) getBalanceFromCache(userID uint) *WalletWithBalance {
	raw, err := svc.redis.Get(getBalanceKey(userID)).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil
		}
		svc.lg.Error("error getting balance from redis", err)
		return nil
	}

	var balance *WalletWithBalance
	if err = json.Unmarshal([]byte(raw), balance); err != nil {
		svc.lg.Error("error unmarshaling balance", err)
		return nil
	}

	return balance
}

func (svc *ServiceImpl) putBalanceToCache(userID uint, balance *WalletWithBalance) {
	raw, err := json.Marshal(balance)
	if err != nil {
		svc.lg.Error("marshal balance to json", err)
	}

	if err = svc.redis.Set(getBalanceKey(userID), string(raw), time.Second*30).Err(); err != nil {
		svc.lg.Error("redis set balance", err)
	}
}
