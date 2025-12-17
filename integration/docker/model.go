package docker

import (
	"github.com/docker/docker/api/types/container"
)

type containerMetadata struct {
	container  *container.Summary
	id         string
	name       string
	protocol   string
	qualifier  qualifier
	portNumber int
}

type qualifier string

const (
	hostQualifier      qualifier = "host"
	containerQualifier qualifier = "container"
)
