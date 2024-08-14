package env

import (
	"errors"
	"fmt"
	"os"
)

type EnvKey string

func (e EnvKey) Require() error {
	if _, ok := os.LookupEnv(string(e)); !ok {
		return fmt.Errorf("env key require, key: %v", string(e))
	}
	return nil
}

func (e EnvKey) Value() string {
	return os.Getenv(string(e))
}

func CheckRequiredEnvs(envs []EnvKey) error {
	var err error
	for _, key := range envs {
		err = errors.Join(err, key.Require())
	}
	return err
}

const (
	NatsURL      EnvKey = "NATS_URL"
	NatsConnName EnvKey = "NATS_CONN_NAME"
	DEV          EnvKey = "DEV"
	JwtSecret    EnvKey = "JWT_KEY"
	RedisHost    EnvKey = "REDIS_HOST"
	MongoURI     EnvKey = "MONGO_URI"
	StripeSecret EnvKey = "STRIPE_SECRET"
)
