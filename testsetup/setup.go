package testsetup

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/matthxwpavin/ticketing/database/mongo"
	"github.com/matthxwpavin/ticketing/env"
	"github.com/matthxwpavin/ticketing/logging/sugar"
	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"
)

func Setup(m *testing.M, conf *mongo.DbConfig, initializer Initializer) {
	// call flag.Parse() here if TestMain uses flags

	logger, err := sugar.New()
	if err != nil {
		log.Fatalf("unable to new logger: %v", err)
	}
	ctx := sugar.WithContext(context.Background(), logger)

	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		logger.Fatalf("Could not construct pool: %s", err)
	}

	// uses pool to try to connect to Docker
	err = pool.Client.Ping()
	if err != nil {
		logger.Fatalf("Could not connect to Docker: %s", err)
	}

	// pulls an image, creates a container based on it and runs it
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "mongo",
		Tag:        "latest",
		Env:        []string{},
	}, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{
			Name: "no",
		}
	})
	if err != nil {
		logger.Fatalf("Could not start resource: %s", err)
	}

	code, testErr := runTests(ctx, m, pool, resource, conf, initializer)
	// You can't defer this because os.Exit doesn't care for defer
	if err := pool.Purge(resource); err != nil {
		logger.Errorw("could not purge resource", "error", err)
	}
	if testErr != nil {
		os.Exit(1)
	}
	os.Exit(code)
}

type Initializer func(context.Context, *mongo.DB) error

func runTests(
	ctx context.Context,
	m *testing.M,
	pool *dockertest.Pool,
	resource *dockertest.Resource,
	conf *mongo.DbConfig,
	initializer Initializer,
) (int, error) {
	logger := sugar.FromContext(ctx)

	os.Setenv("JWT_KEY", "abcd")
	os.Setenv("MONGO_URI", fmt.Sprintf("mongodb://localhost:%s", resource.GetPort("27017/tcp")))
	os.Setenv("NATS_URL", "nats://localhost:4222")
	os.Setenv("NATS_CONN_NAME", "some_name")
	os.Setenv("DEV", "dev")
	if err := env.Load(); err != nil {
		logger.Errorw("unable to load ENV", "error", err)
		return 0, err
	}

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	var db *mongo.DB
	if err := pool.Retry(func() error {
		connCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
		defer cancel()

		db = &mongo.DB{
			URI:    env.MongoURI(),
			Config: *conf,
		}
		if err := db.Connect(connCtx); err != nil {
			logger.Errorw("unable to load database: %v", err)
			return err
		}
		if err := db.Migrate(connCtx); err != nil {
			logger.Errorw("could not migrate the database", "error", err)
			return err
		}
		return nil

	}); err != nil {
		logger.Errorw("could not connect to database", "error", err)
		return 0, err
	}
	// Disconnect the database when main function returns.
	defer db.Disconnect(ctx)

	if err := initializer(ctx, db); err != nil {
		logger.Errorw("could not initialize", "error", err)
		return 0, err
	}

	// Start main test.
	return m.Run(), nil
}
