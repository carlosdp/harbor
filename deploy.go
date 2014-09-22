package main

import (
	"crypto/tls"
	"github.com/coreos/fleet/client"
	"github.com/coreos/fleet/job"
	"github.com/coreos/fleet/registry"
	"github.com/coreos/fleet/ssh"
	"net"
	"net/http"
)

type Deploy struct {
	sshUser     string
	sshHostPort string
}

func NewDeploy(sshUser, sshHostPort string) (*Deploy, error) {
	deploy := &Deploy{
		sshUser:     sshUser,
		sshHostPort: sshHostPort,
	}
	return deploy, nil
}

func (d *Deploy) List() ([]job.Job, error) {
	sshClient, err := ssh.NewSSHClient(d.sshUser, d.sshHostPort, nil, false)

	if err != nil {
		return nil, err
	}

	dial := func(network, addr string) (net.Conn, error) {
		tcpaddr, err := net.ResolveTCPAddr(network, addr)
		if err != nil {
			return nil, err
		}

		return sshClient.DialTCP(network, nil, tcpaddr)
	}

	trans := http.Transport{
		Dial: dial,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	regClient, err := client.NewRegistryClient(&trans, "http://127.0.0.1:4001", registry.DefaultKeyPrefix)

	if err != nil {
		return nil, err
	}

	jobs, err := regClient.Jobs()

	if err != nil {
		return nil, err
	}

	return jobs, nil
}
