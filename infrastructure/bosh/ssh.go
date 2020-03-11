package bosh

import (
	"bytes"
	"github.com/cloudfoundry/bosh-cli/director"
	"github.com/cloudfoundry/bosh-utils/uuid"
	"golang.org/x/crypto/ssh"
	"net"
	"os"
)

// Create a ssh session to a VM
// @param vmId Id of the VM you want to ssh to
// @return session Initialized sessions struct of golang.org/x/crypto/ssh
// @return client Initialized client struct of golang.org/x/crypto/client
func createSshSession(vmId string) (session *ssh.Session, client *ssh.Client, err error){
	uuidGen := uuid.NewGenerator()

	sshOpts, key, err := director.NewSSHOpts(uuidGen)

	if err != nil {
		return nil, nil, err
	}

	// Setup ssh for a specific vm
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

	// Setup config for ssh connection
	config := &ssh.ClientConfig{
		User: entry.Username,
		Auth: authMethods,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// Sets unset values to default
	config.SetDefaults()

	// Create ssh client
	client, err = ssh.Dial("tcp", net.JoinHostPort(entry.Host, "22"), config)
	if err != nil {
		return nil, nil, err
	}

	// Create ssh session with client
	session, err = client.NewSession()
	if err != nil {
		return nil, nil, err
	}

	return session, client, nil
}

// Run a ssh command on a VM
// @param vmId Id of the VM you want to run the command on
// @param command Command you want to execute on the VM
// @return string Stdout of the command execution
func RunSshCommand(vmId string, command string) (string, error) {

	// Create ssh session and client
	session, client, err := createSshSession(vmId)
	defer client.Close()
	defer session.Close()

	if err != nil {
		logError(err, "Failed to create SSH session")
		return "", err
	}

	// Stream session stdout and stderror
	var result bytes.Buffer
	session.Stdout = &result
	session.Stderr = os.Stderr

	// Run the actual command
	err = session.Run(command)

	return result.String(), err
}