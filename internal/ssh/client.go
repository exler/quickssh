package internal_ssh

import (
	"io"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

func NewSSHClient(host, user string, port int, auth Auth) (*ssh.Client, error) {
	client, err := ssh.Dial("tcp", net.JoinHostPort(host, strconv.Itoa(port)), &ssh.ClientConfig{
		User:            user,
		Auth:            auth,
		Timeout:         10 * time.Second,
		HostKeyCallback: VerifyHost,
	})
	if err != nil {
		return nil, err
	}
	return client, err
}

func NewSFTPClient(c *ssh.Client) (*sftp.Client, error) {
	return sftp.NewClient(c)
}

func Run(command string, c *ssh.Client) error {
	session, err := c.NewSession()
	if err != nil {
		return err
	}

	session.Stdout = os.Stdout
	session.Stderr = os.Stderr

	session.Run(command)

	return nil
}

func Download(localPath, remotePath string, c *ssh.Client) (err error) {
	localFile, err := os.Create(localPath)
	if err != nil {
		return
	}
	defer localFile.Close()

	sftp, err := NewSFTPClient(c)
	if err != nil {
		return
	}
	defer sftp.Close()

	remoteFile, err := sftp.Open(remotePath)
	if err != nil {
		return
	}
	defer remoteFile.Close()

	if _, err = io.Copy(localFile, remoteFile); err != nil {
		return
	}

	return localFile.Sync()
}

func Upload(localPath, remotePath string, c *ssh.Client) (err error) {
	localFile, err := os.Open(localPath)
	if err != nil {
		return
	}
	defer localFile.Close()

	sftp, err := NewSFTPClient(c)
	if err != nil {
		return
	}
	defer sftp.Close()

	remoteFile, err := sftp.Create(remotePath)
	if err != nil {
		return
	}
	defer remoteFile.Close()

	_, err = io.Copy(remoteFile, localFile)
	return
}
