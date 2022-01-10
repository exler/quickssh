package cmd

import (
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "quickssh",
	Short: "SSH/SFTP profile manager and client",
	Long: `QuickSSH allows for easy management of SSH profiles and simplifies
working with the SSH and SFTP protocols.`,
}

func Execute() error {
	return Cmd.Execute()
}

func init() {
	Cmd.AddCommand(versionCmd, profileCmd)
}
