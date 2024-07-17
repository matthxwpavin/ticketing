package env

import (
	"os"
	"sync"

	"github.com/matthxwpavin/ticketing/fmts"
)

var (
	jwtkey   string
	mongoURI string
)

var once sync.Once

// Load loads environment variables. It panics if any
// required one is missing.
func Load() {
	once.Do(func() {
		if jwtkey = JWTKey(); jwtkey == "" {
			panicEnv("JWT_KEY")
		}
		if mongoURI = MongoURI(); mongoURI == "" {
			panicEnv("MONGO_URI")
		}

	})
}

func panicEnv(envName string) {
	fmts.Panicf("'%s', env. is required", envName)
}

func JWTKey() string {
	return os.Getenv("JWT_KEY")
}

func MongoURI() string {
	return os.Getenv("MONGO_URI")
}
