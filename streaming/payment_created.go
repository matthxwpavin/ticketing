package streaming

type PaymentCreatedPublisher interface {
	Publisher[PaymentCreatedMessage]
}

type PaymentCreatedConsumer interface {
	JsonConsumer[PaymentCreatedMessage]
}

type PaymentCreatedMessage struct {
}
