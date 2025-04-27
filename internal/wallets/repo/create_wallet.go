package repo

import (
	"context"
	"fmt"

	"github.com/knstch/subtrack-libs/svcerrs"
)

func (r *DBRepo) CreateWallet(ctx context.Context, userID uint, pubKey string, privateKey []byte) error {
	if err := r.db.WithContext(ctx).Create(&Wallet{
		UserID:     userID,
		PublicKey:  pubKey,
		PrivateKey: privateKey,
	}).Error; err != nil {
		if isUniqueViolation(err) {
			return fmt.Errorf("db.Create: %w", svcerrs.ErrConflict)
		}

		return fmt.Errorf("db.Create: %w", err)
	}

	return nil
}
