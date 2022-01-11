package internal_ssh

import (
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/exler/quickssh/internal"
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

func Shell(c *ssh.Client) (err error) {
	session, err := c.NewSession()
	if err != nil {
		return
	}
	defer session.Close()

	stdin, err := session.StdinPipe()
	if err != nil {
		return err
	}
	go io.Copy(stdin, os.Stdin)

	stdout, err := session.StdoutPipe()
	if err != nil {
		return err
	}
	go io.Copy(os.Stdout, stdout)

	stderr, err := session.StderrPipe()
	if err != nil {
		return err
	}
	go io.Copy(os.Stderr, stderr)

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	termWidth, termHeight := internal.GetTerminalSize()
	if err = session.RequestPty("xterm", termWidth, termHeight, modes); err != nil {
		return
	}

	if err = session.Shell(); err != nil {
		return
	}

	signalc := make(chan os.Signal)
	defer func() {
		signal.Reset()
		close(signalc)
	}()
	go propagateSignals(signalc, session, stdin)
	signal.Notify(signalc, os.Interrupt, syscall.SIGTERM)
	return session.Wait()
}

func propagateSignals(signalc chan os.Signal, session *ssh.Session, stdin io.WriteCloser) {
	for s := range signalc {
		switch s {
		case os.Interrupt:
			fmt.Fprint(stdin, "\x03")
		}
	}
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
