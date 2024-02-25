package host_client

import (
	common "github.com/matehaxor03/holistic_common/common"
	"fmt"
)

type HostUser struct {
	GetHost func() Host
	GetUser func() User
	GetFullyQualifiedUsername func() string
	GenerateSSHKey func(other HostUser) []error
}

func newHostUser(host Host, user User) HostUser {
	bashCommand := common.NewBashCommand()
	this_host := host
	this_user := user

	getHost := func() Host {
		return this_host
	}

	getUser := func() User {
		return this_user
	}

	getFullyQualifiedUsername := func() string {
		return getUser().GetUsername() + "@" + getHost().GetHostName()
	}

	generate_ssh_key := func(other HostUser) []error {
		var errors []error
		temp_user := getUser()
		absolute_directory, absolute_directory_errors := temp_user.GetHomeDirectoryAbsoluteDirectory()
		if absolute_directory_errors != nil {
			fmt.Println("absolute_directory_errors")
			return absolute_directory_errors
		}

		absolute_ssh_directory_path := absolute_directory.GetPath()
		absolute_ssh_directory_path = append(absolute_ssh_directory_path, ".ssh")

		ssh_directory, ssh_directory_errors := newAbsoluteDirectory(absolute_ssh_directory_path)
		if ssh_directory_errors != nil {
			fmt.Println("ssh_directory_errors")
			return ssh_directory_errors
		}

		if !ssh_directory.Exists() {
			ssh_directory_create_errors := ssh_directory.Create()
			if ssh_directory_create_errors != nil {
				fmt.Println("ssh_directory_create_errors")
				return ssh_directory_create_errors
			}

			current_user := getUser()
			current_user_group, current_user_group_errors := current_user.GetPrimaryGroup()
			if current_user_group_errors != nil {
				fmt.Println("current_user_group_errors")
				return current_user_group_errors
			}

			if current_user_group == nil {
				errors = append(errors, fmt.Errorf("current user does not have a group"))
				return errors
			}

			set_current_owner_errors := ssh_directory.SetOwnerRecursive(current_user, *current_user_group)
			if set_current_owner_errors != nil {
				fmt.Println("set_current_owner_errors")
				return set_current_owner_errors
			}
		}

		//todo create other user ssh directory
		//todo append ssh pub key to other user if it doesn't exist

		shell_command := "ssh-keygen -b 2048 -t rsa  -f " + ssh_directory.GetPathAsString() + "/" + other.GetFullyQualifiedUsername() + " -C " + other.GetFullyQualifiedUsername() + " -P \"\""
		fmt.Println(shell_command)
		return bashCommand.ExecuteUnsafeCommandSimple(shell_command)
	}

	x := HostUser{
		GetHost: func() Host {
			return getHost()
		},
		GetUser: func() User {
			return getUser()
		},
		GenerateSSHKey: func(other HostUser) []error {
			return generate_ssh_key(other)
		},
		GetFullyQualifiedUsername: func() string {
			return getFullyQualifiedUsername()
		},
	}

	return x
}

