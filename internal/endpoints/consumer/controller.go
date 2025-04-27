package consumer

import (
	"github.com/knstch/wallets-api/event"

	"wallets-service/internal/wallets"

	kafkaPkg "github.com/knstch/subtrack-kafka/consumer"
	"github.com/knstch/subtrack-kafka/topics"
	"github.com/knstch/subtrack-libs/log"
)

type Controller struct {
	lg *log.Logger

	svc wallets.Wallets
}

func NewController(lg *log.Logger, svc wallets.Wallets) *Controller {
	return &Controller{
		lg:  lg,
		svc: svc,
	}
}

func (c *Controller) InitHandlers(consumer *kafkaPkg.Consumer) {
	consumer.AddHandler(topics.TopicWalletsCreateWallet, kafkaPkg.JSONHandler[event.CreateWallet](c.CreateWallet))
}
