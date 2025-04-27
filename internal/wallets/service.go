package wallets

import (
	"context"

	"github.com/knstch/subtrack-libs/log"

	"wallets-service/config"
	"wallets-service/internal/wallets/connections"
	"wallets-service/internal/wallets/repo"
)

type ServiceImpl struct {
	lg *log.Logger

	repo repo.Repository

	conns *connections.Connections

	walletSecret string
}

type Wallets interface {
	CreateWallet(ctx context.Context, userID uint) error
}

func NewService(lg *log.Logger, repo repo.Repository, cfg *config.Config, conns *connections.Connections) *ServiceImpl {

	return &ServiceImpl{
		lg:           lg,
		conns:        conns,
		repo:         repo,
		walletSecret: cfg.WalletSecret,
	}
}
