package consumer

import (
	"context"
	"fmt"
	"github.com/knstch/subtrack-libs/tracing"

	"github.com/knstch/wallets-api/event"
)

func (c *Controller) CreateWallet(ctx context.Context, event *event.CreateWallet) error {
	ctx, span := tracing.StartSpan(ctx, "consumer: CreateWallet")
	defer span.End()

	if err := c.svc.CreateWallet(ctx, event.UserID); err != nil {
		return fmt.Errorf("svc.CreateWallet: %w", err)
	}

	return nil
}
