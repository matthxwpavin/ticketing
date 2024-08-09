package streaming

type ExpirationCompletedPublisher interface {
	Publisher[ExpirationCompletedMessage]
}

type ExpirationCompletedConsumer interface {
	JsonConsumer[ExpirationCompletedMessage]
}

type ExpirationCompletedMessage struct {
	OrderId string `json:"orderId"`
}
