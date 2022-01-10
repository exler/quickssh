package internal_ssh

import (
	"fmt"
	"net"
	"os"
	"syscall"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
	"golang.org/x/term"
)

type Auth []ssh.AuthMethod

func GetAuth(password string) (auth Auth, err error) {
	if HasAgent() {
		auth, err = GetAgent()
		if err != nil {
			return
		}
	} else {
		if password == "" {
			passwordBytes, _ := term.ReadPassword(int(syscall.Stdin))
			password = string(passwordBytes)
		}

		auth = Password(password)
	}
	return
}

func Password(password string) Auth {
	return Auth{
		ssh.Password(password),
	}
}

func Key(privateKeyPath, passphrase string) (Auth, error) {
	signer, err := GetSigner(privateKeyPath, passphrase)
	if err != nil {
		return nil, err
	}

	return Auth{
		ssh.PublicKeys(signer),
	}, nil
}

func HasAgent() bool {
	return os.Getenv("SSH_AUTH_SOCK") != ""
}

func GetAgent() (Auth, error) {
	sshAgent, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK"))
	if err != nil {
		return nil, fmt.Errorf("could not find SSH agent: %w", err)
	}

	return Auth{
		ssh.PublicKeysCallback(agent.NewClient(sshAgent).Signers),
	}, nil
}

func GetSigner(privateKeyPath, passphrase string) (signer ssh.Signer, err error) {
	privateKey, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return nil, err
	}

	if passphrase == "" {
		return ssh.ParsePrivateKey(privateKey)
	} else {
		return ssh.ParsePrivateKeyWithPassphrase(privateKey, []byte(passphrase))
	}
}
