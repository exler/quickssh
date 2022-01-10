package cmd

import (
	"fmt"
	"log"

	"github.com/exler/quickssh/internal"
	internal_ssh "github.com/exler/quickssh/internal/ssh"
	"github.com/spf13/cobra"
)

var (
	fileCmd = &cobra.Command{
		Use:   "file",
		Short: "Download or upload files using SFTP",
	}

	downloadFileCmd = &cobra.Command{
		Use:   "download",
		Short: "Download a file from the remote server",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				fmt.Println("Provide the profile name or target in format 'user@host` as an argument.")
				return
			}

			if len(args) < 2 {
				fmt.Println("Provide the path to the file you want to download.")
				return
			}

			if len(args) < 3 {
				fmt.Println("Provide the path where the file should be downloaded.")
				return
			}

			target := args[0]
			remotePath := args[1]
			localPath := args[2]

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

			err = internal_ssh.Download(localPath, remotePath, client)
			if err != nil {
				log.Printf("Download error: %s", err)
				return
			}
		},
	}

	uploadFileCmd = &cobra.Command{
		Use:   "upload",
		Short: "Upload a file to the remote server",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				fmt.Println("Provide the profile name or target in format 'user@host` as an argument.")
				return
			}

			if len(args) < 2 {
				fmt.Println("Provide the path to the file you want to upload.")
				return
			}

			if len(args) < 3 {
				fmt.Println("Provide the path where the file should be uploaded.")
				return
			}

			target := args[0]
			localPath := args[1]
			remotePath := args[2]

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

			err = internal_ssh.Upload(localPath, remotePath, client)
			if err != nil {
				log.Printf("Upload error: %s", err)
				return
			}
		},
	}
)

func init() {
	fileCmd.AddCommand(downloadFileCmd, uploadFileCmd)
}
