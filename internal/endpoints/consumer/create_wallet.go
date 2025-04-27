package consumer

import (
	"context"
	"fmt"

	"github.com/knstch/wallets-api/event"
)

func (c *Controller) CreateWallet(ctx context.Context, event *event.CreateWallet) error {
	if err := c.svc.CreateWallet(ctx, event.UserID); err != nil {
		return fmt.Errorf("svc.CreateWallet: %w", err)
	}

	return nil
}
