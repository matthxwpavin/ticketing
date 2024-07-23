package env

import (
	"os"
	"sync"

	"github.com/matthxwpavin/ticketing/fmts"
)

type Env struct {
	jwtKey       string
	mongoURI     string
	natsURL      string
	natsConnName string
}

var env Env

var once sync.Once

// Load loads environment variables. It panics if any
// required one is missing.
func Load() {
	once.Do(func() {

		envMap := map[string]*string{
			"JWT_KEY":        &env.jwtKey,
			"MONGO_URI":      &env.mongoURI,
			"NATS_URL":       &env.natsURL,
			"NATS_CONN_NAME": &env.natsConnName,
		}
		for k, dst := range envMap {
			load(k, dst)
		}
	})
}

func load(key string, dst *string) {
	if *dst = os.Getenv(key); *dst == "" {
		panicEnv(key)
	}
}

func panicEnv(envName string) {
	fmts.Panicf("'%s', env. is required", envName)
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
