package integrationtest

import (
	"log"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/Joyang0419/backendmodules/database/mysql/client"
	"github.com/ory/dockertest"
)

func CreateContainer(name string) (*dockertest.Pool, *dockertest.Resource, *gorm.DB) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	// pulls an image, creates a container based on it and runs it
	options := &dockertest.RunOptions{
		Name:       name,
		Repository: "mysql",
		Tag:        "latest",
		Env: []string{
			"MYSQL_ROOT_PASSWORD=root",
			"MYSQL_DATABASE=dev",
			"MYSQL_USER=joy",
			"MYSQL_PASSWORD=joy",
		},
	}
	resource, err := pool.RunWithOptions(options)
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)

	}

	host, port := GetHostPort(resource, "3306/tcp")
	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	var dbConn *gorm.DB
	if err = pool.Retry(func() error {
		var retryErr error
		if dbConn, retryErr = client.SetupDB(
			client.Config{
				Host:            host,
				Port:            port,
				Username:        "joy",
				Password:        "joy",
				Database:        "dev",
				MaxIdleConns:    10,
				MaxOpenConns:    10,
				ConnMaxLifeTime: 10 * time.Minute,
			}, logger.Silent,
		); retryErr != nil {
			return retryErr
		}
		return nil

	}); err != nil {
		// You can't defer this because os.Exit doesn't care for defer
		if errPurge := pool.Purge(resource); errPurge != nil {
			log.Fatalf("Could not purge resource: %s", err)
		}
	}

	return pool, resource, dbConn
}
