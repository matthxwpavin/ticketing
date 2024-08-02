package env

import (
	"fmt"
	"os"
	"sync"
)

type Env struct {
	jwtKey       string
	mongoURI     string
	natsURL      string
	natsConnName string
	dev          string
}

var env Env

var envMap = map[string]*string{
	"JWT_KEY":        &env.jwtKey,
	"MONGO_URI":      &env.mongoURI,
	"NATS_URL":       &env.natsURL,
	"NATS_CONN_NAME": &env.natsConnName,
	"DEV":            &env.dev,
}

var once sync.Once

// Load loads environment variables. It returns error if any
// required one is missed.
func Load() (err error) {
	once.Do(func() {
		for k, dst := range envMap {
			if err = load(k, dst); err != nil {
				break
			}
		}
	})
	return
}

func load(key string, dst *string) error {
	if *dst = os.Getenv(key); *dst == "" {
		return fmtError(key)
	}
	return nil
}

func fmtError(envName string) error {
	return fmt.Errorf("'%s', env. is required", envName)
}

func JWTKey() string {
	return env.jwtKey
}

func MongoURI() string {
	return env.mongoURI
}

func NatsURL() string {
	return env.natsURL
}

func NatsConnectionName() string {
	return env.natsConnName
}

func Dev() bool {
	return env.dev == "dev"
}
