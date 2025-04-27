package repo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/knstch/subtrack-libs/svcerrs"
	"gorm.io/gorm"

	"wallets-service/internal/domain/dto"
)

func (r *DBRepo) GetWalletByUserID(ctx context.Context, userID uint) (dto.Wallet, error) {
	tx := r.db.WithContext(ctx)

	var wallet Wallet
	if err := tx.Preload(UserTokens{}.TableName()).Where("user_id = ?", userID).First(&wallet).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.Wallet{}, fmt.Errorf("wallet not found: %w", svcerrs.ErrDataNotFound)
		}
		return dto.Wallet{}, fmt.Errorf("db.First: %w", err)
	}

	var userTokens UserTokens
	if err := tx.Where("wallet_id = ?", wallet.ID).First(&userTokens).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.Wallet{
				PublicAddr: wallet.PublicKey,
			}, nil
		}
		return dto.Wallet{}, fmt.Errorf("db.First: %w", err)
	}

	tokens := make([]string, 0)
	if err := json.Unmarshal(userTokens.Tokens, &tokens); err != nil {
		return dto.Wallet{}, fmt.Errorf("json.Unmarshal: %w", err)
	}

	return dto.Wallet{
		PublicAddr:     wallet.PublicKey,
		TokenAddresses: tokens,
	}, nil
}
