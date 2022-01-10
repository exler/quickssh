package cmd

import (
	"fmt"
	"strconv"
	"syscall"

	internal "github.com/exler/quickssh/internal"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var (
	profileCmd = &cobra.Command{
		Use:   "profile",
		Short: "Manage server profiles",
	}

	listProfileCmd = &cobra.Command{
		Use:   "list",
		Short: "List available profiles",
		Run: func(cmd *cobra.Command, args []string) {
			profiles := internal.GetProfiles()
			profileCnt := 1
			for profileName, profile := range profiles {
				fmt.Printf("%d) %s\n\tHostname: %s\n\tPort: %d\n\tUsername: %s\n", profileCnt, profileName, profile.Hostname, profile.Port, profile.Username)
				profileCnt++
			}
		},
	}

	addProfileCmd = &cobra.Command{
		Use:   "add",
		Short: "Add new profile",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 1 {
				fmt.Println("Provide the profile name as an argument.")
				return
			}
			profileName := args[0]

			profiles := internal.GetProfiles()
			if _, ok := profiles[profileName]; ok {
				fmt.Println("Profile already exists!")
				return
			}

			fmt.Print("Hostname: ")
			host := internal.GetUserInput()

			fmt.Print("Port (default: 22): ")
			var port int
			portStr := internal.GetUserInput()
			if portStr == "" {
				port = 22
			} else {
				port, _ = strconv.Atoi(portStr)
			}

			fmt.Print("Username: ")
			user := internal.GetUserInput()

			fmt.Print("Password (optional): ")
			pass, _ := term.ReadPassword(int(syscall.Stdin))

			profile := internal.Profile{
				Hostname: host,
				Port:     port,
				Username: user,
				Password: string(pass),
			}

			profiles[profileName] = profile

			internal.SetProfiles(profiles)
			fmt.Println("\nProfile saved successfully!")
		},
	}

	editProfileCmd = &cobra.Command{
		Use:   "edit",
		Short: "Edit existing profile",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 1 {
				fmt.Println("Provide the profile name as an argument.")
				return
			}
			profileName := args[0]

			profiles := internal.GetProfiles()
			var profile internal.Profile
			var ok bool
			if profile, ok = profiles[profileName]; !ok {
				fmt.Println("No profile with that name exists!")
				return
			}

			fmt.Printf("Hostname (currently: %s): ", profile.Hostname)
			host := internal.GetUserInput()
			if host == "" {
				host = profile.Hostname
			}

			fmt.Printf("Port (currently: %d): ", profile.Port)
			var port int
			portStr := internal.GetUserInput()
			if portStr == "" {
				port = profile.Port
			} else {
				port, _ = strconv.Atoi(portStr)
			}

			fmt.Printf("Username (currently: %s): ", profile.Username)
			user := internal.GetUserInput()
			if user == "" {
				user = profile.Username
			}

			fmt.Print("Password (optional): ")
			pass, _ := term.ReadPassword(int(syscall.Stdin))

			changedProfile := internal.Profile{
				Hostname: host,
				Port:     port,
				Username: user,
				Password: string(pass),
			}

			profiles[profileName] = changedProfile

			internal.SetProfiles(profiles)
			fmt.Println("\nProfile saved successfully!")
		},
	}

	deleteProfileCmd = &cobra.Command{
		Use:   "delete",
		Short: "Delete existing profile",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 1 {
				fmt.Println("Provide the profile name as an argument.")
				return
			}
			profileName := args[0]

			profiles := internal.GetProfiles()
			if _, ok := profiles[profileName]; ok {
				delete(profiles, profileName)
				internal.SetProfiles(profiles)
				fmt.Println("Profile deleted!")
			} else {
				fmt.Println("No profile with that name exists!")
			}
		},
	}
)

func init() {
	profileCmd.AddCommand(listProfileCmd, addProfileCmd, editProfileCmd, deleteProfileCmd)
}
