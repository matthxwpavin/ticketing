package env

import (
	"os"
	"sync"

	"github.com/matthxwpavin/ticketing/fmts"
)

type Env struct {
	jwtKey   string
	mongoURI string
	natsURL  string
}

var env *Env

var once sync.Once

// Load loads environment variables. It panics if any
// required one is missing.
func Load() {
	once.Do(func() {
		env := &Env{}

		for _, load := range []func(*Env){
			loadJWTKey,
			loadMongoURI,
			loadNatsURL,
		} {
			load(env)
		}
	})
}

func loadJWTKey(env *Env) {
	const key = "JWT_KEY"
	if env.jwtKey = os.Getenv(key); env.jwtKey == "" {
		panicEnv(key)
	}
}

func loadMongoURI(env *Env) {
	const key = "MONGO_URI"
	if env.mongoURI = os.Getenv(key); env.mongoURI == "" {
		panicEnv(key)
	}
}

func loadNatsURL(env *Env) {
	const key = "NATS_URL"
	if env.natsURL = os.Getenv(key); env.natsURL == "" {
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
