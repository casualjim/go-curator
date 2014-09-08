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
	zkURL       = "^(?:(?:zk|zookeeper)://)?(" + hostAndPort + "(?:," + hostAndPort + ")*)/?(:?(" + zNode + "(?:/" + zNode + ")*))*$"
)

var (
	zkURLRegex = regexp.MustCompile(zkURL)
)

// Provider provides the zookeeper connection string
type Provider interface {
	io.Closer
	Start() error
	Hosts() []string
}

type fixedEnsembleProvider struct {
	hosts []string
}

// Start starts this ensemble provider
func (p *fixedEnsembleProvider) Start() error {
	return nil
}

// Close stops this ensemble provider
func (p *fixedEnsembleProvider) Close() error {
	return nil
}

// Hosts returns list of zookeeper hosts with their ports [name:port] for this ensemble
func (p *fixedEnsembleProvider) Hosts() []string {
	return p.hosts
}

// Fixed returns a fixed ensemble provider, fixed to the ensemble provded by the uri string
func Fixed(uri string) (Provider, string, error) {
	hosts, path, err := parseZookeeperURI(uri)
	return &fixedEnsembleProvider{hosts: hosts}, path, err
}

func parseZookeeperURI(uri string) ([]string, string, error) {
	parsed := zkURLRegex.FindStringSubmatch(uri)
	if len(parsed) < 2 {
		return nil, "", errors.New("The uri " + uri + " is not a valid connection string for zookeeper")
	}
	pth := ""
	if len(parsed) > 2 {
		pth = parsed[2]
	}
	return strings.Split(parsed[1], ","), "/" + pth, nil
}
