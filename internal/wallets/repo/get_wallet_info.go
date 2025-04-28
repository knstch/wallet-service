package repo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/knstch/subtrack-libs/svcerrs"
	"gorm.io/gorm"
)

func (r *DBRepo) GetWalletInfoByUserID(ctx context.Context, userID uint) (WalletInfo, error) {
	tx := r.db.WithContext(ctx)

	var wallet Wallet
	if err := tx.Where("user_id = ?", userID).First(&wallet).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return WalletInfo{}, fmt.Errorf("wallet not found: %w", svcerrs.ErrDataNotFound)
		}
		return WalletInfo{}, fmt.Errorf("db.First: %w", err)
	}

	var userTokens UserTokens
	if err := tx.Where("wallet_id = ?", wallet.ID).First(&userTokens).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return WalletInfo{
				PublicAddr: wallet.PublicKey,
			}, nil
		}
		return WalletInfo{}, fmt.Errorf("db.First: %w", err)
	}

	tokens := make([]string, 0)
	if err := json.Unmarshal(userTokens.Tokens, &tokens); err != nil {
		return WalletInfo{}, fmt.Errorf("json.Unmarshal: %w", err)
	}

	return WalletInfo{
		PublicAddr:     wallet.PublicKey,
		TokenAddresses: tokens,
	}, nil
}

func (r *DBRepo) GetWallet(ctx context.Context, userID uint) (string, []byte, error) {
	var wallet Wallet
	if err := r.db.WithContext(ctx).First(&wallet, "user_id = ?", userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", nil, fmt.Errorf("db.First: %w", svcerrs.ErrDataNotFound)
		}
		return "", nil, fmt.Errorf("db.First: %w", err)
	}

	return wallet.PublicKey, wallet.PrivateKey, nil
}
