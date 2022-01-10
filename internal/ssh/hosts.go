package internal_ssh

import (
	"errors"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"

	"github.com/exler/quickssh/internal"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
)

func DefaultKnownHostsPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, ".ssh", "known_hosts"), err
}

func DefaultKnownHosts() (ssh.HostKeyCallback, error) {
	path, err := DefaultKnownHostsPath()
	if err != nil {
		return nil, err
	}

	return knownhosts.New(path)
}

func CheckKnownHost(host string, remote net.Addr, key ssh.PublicKey, knownHostsFile string) (found bool, err error) {
	var keyErr *knownhosts.KeyError

	if knownHostsFile == "" {
		path, err := DefaultKnownHostsPath()
		if err != nil {
			return false, err
		}

		knownHostsFile = path
	}

	callback, err := knownhosts.New(knownHostsFile)
	if err != nil {
		return false, err
	}

	err = callback(host, remote, key)

	// Known host exists
	if err == nil {
		return true, nil
	}

	// Make sure that the error returned from the callback is host not in file error.
	// If keyErr.Want is greater than 0 length, that means host is in file with different key.
	if errors.As(err, &keyErr) && len(keyErr.Want) > 0 {
		return true, keyErr
	}

	// Some other error occurred and safest way to handle is to pass it back to user.
	if err != nil {
		return false, err
	}

	// Key is not trusted because it is not in the file.
	return false, nil
}

func AddKnownHost(host string, remote net.Addr, key ssh.PublicKey, knownHostsFile string) (err error) {
	if knownHostsFile == "" {
		path, err := DefaultKnownHostsPath()
		if err != nil {
			return err
		}

		knownHostsFile = path
	}

	f, err := os.OpenFile(knownHostsFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		return err
	}

	defer f.Close()

	remoteNormalized := knownhosts.Normalize(remote.String())
	hostNormalized := knownhosts.Normalize(host)
	addresses := []string{remoteNormalized}

	if hostNormalized != remoteNormalized {
		addresses = append(addresses, hostNormalized)
	}

	_, err = f.WriteString(knownhosts.Line(addresses, key) + "\n")
	return err
}

func VerifyHost(host string, remote net.Addr, key ssh.PublicKey) error {
	// Check if host is in the known_hosts file
	hostFound, err := CheckKnownHost(host, remote, key, "")

	// Host is in known_hosts but the keys are mismatched
	if hostFound {
		if err != nil {
			return err
		} else {
			return nil
		}
	}

	// Ask the user whether he trusts the host
	if !askIsHostTrusted(host, key) {
		return errors.New("host key untrusted")
	}

	// Add the new host to known hosts file
	return AddKnownHost(host, remote, key, "")
}

func askIsHostTrusted(host string, key ssh.PublicKey) bool {
	fmt.Printf("Unknown Host: %s \nFingerprint: %s \n", host, ssh.FingerprintSHA256(key))
	fmt.Print("Would you like to add it? [Y/n]: ")

	answer := internal.GetUserInput()
	if strings.ToLower(answer) == "n" {
		return false
	} else {
		return true
	}
}
