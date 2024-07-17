package env

import (
	"os"
	"sync"

	"github.com/matthxwpavin/ticketing/fmts"
)

var (
	jwtkey string
)

var once sync.Once

// Load loads environment variables. It panics if any
// required one is missing.
func Load() {
	once.Do(func() {
		if jwtkey = GetJWTKey(); jwtkey == "" {
			fmts.Panicf("'JWT_KEY' env is required")
		}
	})
}

func GetJWTKey() string {
	return os.Getenv("JWT_KEY")
}
