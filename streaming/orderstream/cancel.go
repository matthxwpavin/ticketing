package orderstream

import "context"

type CancelledPublisher interface {
	Publish(context.Context, *CancelledMessage) error
}

type CancelledMessage struct {
	OrderId string `json:"orderId"`
	Ticket  struct {
		Id string `json:"id"`
	} `json:"ticket"`
}
