package bosh

import (
	"github.com/cloudfoundry/bosh-cli/director"
	"github.com/cloudfoundry/bosh-utils/uuid"
	"golang.org/x/crypto/ssh"
	"net"
)

func createSshSession(vmId string) (session *ssh.Session, client *ssh.Client, err error){
	uuidGen := uuid.NewGenerator()

	sshOpts, key, err := director.NewSSHOpts(uuidGen)

	if err != nil {
		return nil, nil, err
	}

	sshResult, err := deployment.SetUpSSH(director.NewAllOrInstanceGroupOrInstanceSlug("", vmId),
		sshOpts)

	if err != nil {
		return nil, nil, err
	}

	entry := sshResult.Hosts[0]

	var authMethods []ssh.AuthMethod

	if key != "" {
		signer, err := ssh.ParsePrivateKey([]byte(key))
		if err != nil {
			return nil, nil, err
		}
		authMethods = append(authMethods, ssh.PublicKeys(signer))
	}

	config := &ssh.ClientConfig{
		User: entry.Username,
		Auth: authMethods,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	config.SetDefaults()

	client, err = ssh.Dial("tcp", net.JoinHostPort(entry.Host, "22"), config)
	if err != nil {
		return nil, nil, err
	}

	session, err = client.NewSession()
	if err != nil {
		return nil, nil, err
	}

	return session, client, nil
}
