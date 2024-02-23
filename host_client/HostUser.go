package host_client

import (
	common "github.com/matehaxor03/holistic_common/common"
)

type HostUser struct {
	GetHost func() Host
	GetUser func() User
	GetFullyQualifiedUsername func() string
	GenerateSSHKey func(other HostUser) []error
}

func newHostUser(host Host, user User) (*HostUser, []error) {
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
		temp_user := getUser()
		absolute_directory, absolute_directory_errors := temp_user.GetHomeDirectoryAbsoluteDirectory()
		if absolute_directory_errors != nil {
			return absolute_directory_errors
		}

		absolute_ssh_directory_path := absolute_directory.GetPath()
		absolute_ssh_directory_path = append(absolute_ssh_directory_path, ".ssh")

		ssh_directory, ssh_directory_errors := newAbsoluteDirectory(absolute_ssh_directory_path)
		if ssh_directory_errors != nil {
			return ssh_directory_errors
		}

		if !ssh_directory.Exists() {
			ssh_directory_create_errors := ssh_directory.Create()
			if ssh_directory_create_errors != nil {
				return ssh_directory_create_errors
			}

			//todo set ownership primary user id and others
		}

		shell_command := "ssh-keygen -b 2048 -t rsa  -f " + ssh_directory.GetPathAsString() + "/" + other.GetFullyQualifiedUsername() + " -C " + other.GetFullyQualifiedUsername() + " -P \"\""
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

	return &x, nil
}

