package integrationtest_test

import (
	"log"
	"testing"

	"github.com/Joyang0419/backendmodules/database/mysql/integrationtest"
	"github.com/stretchr/testify/assert"
)

func TestCreateContainer(t *testing.T) {
	pool, resource, dbConn := integrationtest.CreateContainer("mysqldb")
	defer func() {
		if errPurge := pool.Purge(resource); errPurge != nil {
			log.Fatalf("Could not purge resource: %s", errPurge)
		}
	}()

	assert.Equal(t, "mysql", dbConn.Name())
}
