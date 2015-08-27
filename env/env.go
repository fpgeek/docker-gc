package env

import (
	"strconv"

	log "github.com/Sirupsen/logrus"
	"github.com/frhwang/gopkg/env"
)

// Environments
var (
	DockerHost          = env.GetOrDefault("DOCKER_HOST", "unix:///var/run/docker.sock")
	DockerCertPath      = env.GetOrDefault("DOCKER_CERT_PATH", "")
	KeepStaleImageCount int
)

func init() {
	intValue, err := strconv.ParseInt(env.GetOrDefault("KEEP_STALE_IMAGE_COUNT", "1"), 10, 0)
	if err != nil {
		log.Fatal(err)
	}
	KeepStaleImageCount = int(intValue)
}
