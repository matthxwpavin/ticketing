package testsetup

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/matthxwpavin/ticketing/database/mongo"
	"github.com/matthxwpavin/ticketing/logging/sugar"
	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"
)

func Setup(m *testing.M, conf *mongo.DbConfig, onInit InitFunc) {
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

	// Defers a function to purge the resource when any panicking occurred
	// before Main.Run will be executed. Panicking inside of sub-tests(T)
	// not be recovered by this defer.
	defer func() {
		if a := recover(); a != nil {
			logger.Errorw("a panic recovered", "error", a)
			purge(ctx, pool, resource)
			// Re-throw the panic to exit the program with a non-zero code
			// to mark this test as a failed one.
			panic(a)
		}
	}()

	code, testErr := runTests(ctx, m, pool, resource, conf, onInit)
	// You can't defer this because os.Exit doesn't care for defer
	purge(ctx, pool, resource)

	if testErr != nil {
		os.Exit(1)
	}
	os.Exit(code)
}

type InitFunc func(context.Context, *mongo.DB) error

// purge purges the resources (running container). It usually be called in the last
// step of testing (before os.Exit with a code produced by Main.Run).
func purge(ctx context.Context, pool *dockertest.Pool, resource *dockertest.Resource) {
	logger := sugar.FromContext(ctx)
	if err := pool.Purge(resource); err != nil {
		logger.Errorw("could not purge resource", "error", err)
	} else {
		logger.Infoln("resources purged")
	}
}

func runTests(
	ctx context.Context,
	m *testing.M,
	pool *dockertest.Pool,
	resource *dockertest.Resource,
	conf *mongo.DbConfig,
	onInit InitFunc,
) (int, error) {
	logger := sugar.FromContext(ctx)

	os.Setenv("MONGO_URI", fmt.Sprintf("mongodb://localhost:%s", resource.GetPort("27017/tcp")))

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	var db *mongo.DB
	if err := pool.Retry(func() error {
		connCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
		defer cancel()

		db = &mongo.DB{
			URI:    fmt.Sprintf("mongodb://localhost:%s", resource.GetPort("27017/tcp")),
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

	if onInit != nil {
		if err := onInit(ctx, db); err != nil {
			logger.Errorw("could not initialize", "error", err)
			return 0, err
		}
	}

	// Start main test.
	return m.Run(), nil
}
