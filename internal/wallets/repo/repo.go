package repo

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
)

type Repository interface {
	CreateWallet(ctx context.Context, userID uint, pubKey string, privateKey []byte) error
	GetWalletInfoByUserID(ctx context.Context, userID uint) (WalletInfo, error)
	GetWallet(ctx context.Context, userID uint) (string, []byte, error)

	Transaction(fn func(st Repository) error) error
}

type Wallet struct {
	ID         uint
	UserID     uint
	PublicKey  string
	PrivateKey []byte
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  *time.Time
}

func (Wallet) TableName() string {
	return "wallets"
}

type UserTokens struct {
	ID        uint
	WalletID  uint
	Tokens    []byte
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

func (UserTokens) TableName() string {
	return "user_tokens"
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}
	return false
}
