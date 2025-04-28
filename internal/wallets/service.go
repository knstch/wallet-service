package wallets

import (
	"context"
	"wallets-service/internal/domain/enum"
	"wallets-service/internal/wallets/connections"

	"github.com/knstch/subtrack-libs/log"

	"github.com/go-redis/redis"

	"wallets-service/config"
	"wallets-service/internal/wallets/repo"
)

type ServiceImpl struct {
	lg           *log.Logger
	repo         repo.Repository
	walletSecret string
	redis        *redis.Client
	connections  *connections.Connections
}

type Wallets interface {
	CreateWallet(ctx context.Context, userID uint) error
	GetBalance(ctx context.Context, network enum.Network) (*WalletWithBalance, error)
}

func NewService(lg *log.Logger, repo repo.Repository, cfg *config.Config, connections *connections.Connections, redis *redis.Client) *ServiceImpl {
	return &ServiceImpl{
		lg:           lg,
		repo:         repo,
		walletSecret: cfg.WalletSecret,
		connections:  connections,
		redis:        redis,
	}
}
