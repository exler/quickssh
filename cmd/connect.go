package cmd

import (
	"fmt"
	"log"

	"github.com/exler/quickssh/internal"
	internal_ssh "github.com/exler/quickssh/internal/ssh"
	"github.com/spf13/cobra"
)

var connectPort int
var connectCmd = &cobra.Command{
	Use:   "connect",
	Short: "Connect to a SSH server",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("Provide the profile name or target in format 'user@host` as an argument.")
			return
		}

		target := args[0]

		profiles := internal.GetProfiles()
		var profile internal.Profile
		var ok bool
		if profile, ok = profiles[target]; !ok {
			profile.Username, profile.Hostname = internal.ParseUserAndHost(target)
		}

		auth, err := internal_ssh.GetAuth(profile.Password)
		if err != nil {
			log.Printf("Authentication error: %s", err)
			return
		}

		client, err := internal_ssh.NewSSHClient(profile.Hostname, profile.Username, profile.Port, auth)
		if err != nil {
			log.Printf("SSH connection error: %s", err)
			return
		}
		defer client.Close()
	},
}

func init() {
	connectCmd.PersistentFlags().IntVarP(&connectPort, "port", "p", 22, "SSH port to connect to")
}
