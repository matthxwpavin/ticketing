package streaming

type ChargeCreatedPublisher interface {
	Publisher[ChargeCreatedMessage]
}

type ChargeCreatedConsumer interface {
	JsonConsumer[ChargeCreatedMessage]
}

type ChargeCreatedMessage struct {
}
