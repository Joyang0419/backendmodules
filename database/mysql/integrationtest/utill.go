package integrationtest

import (
	"log"
	"net/url"
	"os"

	"github.com/ory/dockertest"
)

func GetHostPort(resource *dockertest.Resource, id string) (host, port string) {
	dockerURL := os.Getenv("DOCKER_HOST")
	if dockerURL == "" {
		return "localhost", resource.GetPort(id)
	}

	u, err := url.Parse(dockerURL)
	if err != nil {
		log.Fatal(err)
	}

	return u.Hostname(), u.Port()
}
