// Package ensemble provides an interface and a set of packages
// to provide zookeeper ensembles
package ensemble

import (
	"errors"
	"io"
	"regexp"
	"strings"
)

const (
	hostAndPort = "[A-z0-9-.]+(?::\\d+)?"
	zNode       = "[^/]+"
	zkUrl       = "^(?:(?:zk|zookeeper)://)?(" + hostAndPort + "(?:," + hostAndPort + ")*)/?(:?(" + zNode + "(?:/" + zNode + ")*))*$"
)

var (
	zkUrlRegex = regexp.MustCompile(zkUrl)
)

// EnsembleProvider provides the zookeeper connection string
type EnsembleProvider interface {
	io.Closer
	Start() error
	Hosts() []string
}

// FixedEnsembleProvider is a constant ensemble, it never changes and is created with a connection string
type FixedEnsembleProvider struct {
	hosts []string
}

func (p *FixedEnsembleProvider) Start() error {
	return nil
}

func (p *FixedEnsembleProvider) Close() error {
	return nil
}

func (p *FixedEnsembleProvider) Hosts() []string {
	return p.hosts
}

func New(uri string) (EnsembleProvider, string, error) {
	hosts, path, err := parseZookeeperUri(uri)
	return &FixedEnsembleProvider{hosts: hosts}, path, err
}

func parseZookeeperUri(uri string) ([]string, string, error) {
	parsed := zkUrlRegex.FindStringSubmatch(uri)
	if len(parsed) < 2 {
		return nil, "", errors.New("The uri " + uri + " is not a valid connection string for zookeeper")
	}
	pth := ""
	if len(parsed) > 2 {
		pth = parsed[2]
	}
	return strings.Split(parsed[1], ","), "/" + pth, nil
}
