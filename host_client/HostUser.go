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

	create_ssh_directory := func(in_user User) (*AbsoluteDirectory, []error) {
		var errors []error
		absolute_directory, absolute_directory_errors := in_user.GetHomeDirectoryAbsoluteDirectory()
		if absolute_directory_errors != nil {
			fmt.Println("absolute_directory_errors")
			return nil, absolute_directory_errors
		}

		absolute_ssh_directory_path := absolute_directory.GetPath()
		absolute_ssh_directory_path = append(absolute_ssh_directory_path, ".ssh")

		ssh_directory, ssh_directory_errors := newAbsoluteDirectory(absolute_ssh_directory_path)
		if ssh_directory_errors != nil {
			fmt.Println("ssh_directory_errors")
			return nil, ssh_directory_errors
		}

		current_user := getUser()
		current_user_group, current_user_group_errors := current_user.GetPrimaryGroup()
		if current_user_group_errors != nil {
			fmt.Println("current_user_group_errors")
			return nil, current_user_group_errors
		}

		if current_user_group == nil {
			errors = append(errors, fmt.Errorf("current user does not have a group"))
			return nil, errors
		}

		if !ssh_directory.Exists() {
			ssh_directory_create_errors := ssh_directory.Create()
			if ssh_directory_create_errors != nil {
				fmt.Println("ssh_directory_create_errors")
				return nil, ssh_directory_create_errors
			}

			set_current_owner_errors := ssh_directory.SetOwnerRecursive(current_user, *current_user_group)
			if set_current_owner_errors != nil {
				fmt.Println("set_current_owner_errors")
				return nil, set_current_owner_errors
			}
		}

		return ssh_directory, nil
	}

	create_ssh_key := func(source_ssh_directory AbsoluteDirectory, destination_host_user HostUser, destination_ssh_directory AbsoluteDirectory) []error {
		var errors []error
		source_user := getUser()
		destination_user := destination_host_user.GetUser()
		
		source_user_group, source_user_group_errors := source_user.GetPrimaryGroup()
		if source_user_group_errors != nil {
			fmt.Println("source_user_group_errors")
			return source_user_group_errors
		}

		if source_user_group == nil {
			errors = append(errors, fmt.Errorf("source user does not have a group"))
			return errors
		}

		destination_user_group, destination_user_group_errors := destination_user.GetPrimaryGroup()
		if destination_user_group_errors != nil {
			fmt.Println("destination_user_group_errors")
			return destination_user_group_errors
		}

		if destination_user_group == nil {
			errors = append(errors, fmt.Errorf("destination user does not have a group"))
			return errors
		}
		
		source_file_ssh_private_key, source_file_ssh_private_key_errors := newAbsoluteFile(source_ssh_directory, destination_host_user.GetFullyQualifiedUsername())
		if source_file_ssh_private_key_errors != nil {
			fmt.Println("source_file_ssh_private_key_errors")
			return source_file_ssh_private_key_errors
		}

		source_file_ssh_private_key_remove_errors := source_file_ssh_private_key.RemoveIfExists()
		if source_file_ssh_private_key_remove_errors != nil {
			fmt.Println("source_file_ssh_private_key_remove_errors")
			return source_file_ssh_private_key_remove_errors
		}

		
		source_file_ssh_public_key, source_file_ssh_public_key_errors := newAbsoluteFile(source_ssh_directory, destination_host_user.GetFullyQualifiedUsername() + ".pub")
		if source_file_ssh_public_key_errors != nil {
			fmt.Println("source_file_ssh_public_key_errors")
			return source_file_ssh_public_key_errors
		}

		source_file_ssh_public_key_remove_errors := source_file_ssh_public_key.RemoveIfExists()
		if source_file_ssh_public_key_remove_errors != nil {
			fmt.Println("source_file_ssh_public_key_remove_errors")
			return source_file_ssh_public_key_remove_errors
		}
		

		shell_command := "ssh-keygen -b 2048 -t rsa  -f " + source_ssh_directory.GetPathAsString() + "/" + destination_host_user.GetFullyQualifiedUsername() + " -C " + destination_host_user.GetFullyQualifiedUsername() + " -P \"\""
		bash_command_errors := bashCommand.ExecuteUnsafeCommandSimple(shell_command)

		if bash_command_errors != nil {
			fmt.Println("bash_command_errors")
			return bash_command_errors
		}

		set_source_current_owner_absoloute_file_ssh_private_key_errors := source_file_ssh_private_key.SetOwner(source_user, *source_user_group)
		if set_source_current_owner_absoloute_file_ssh_private_key_errors != nil {
			fmt.Println("set_current_owner_absoloute_file_ssh_private_key_errors")
			return set_source_current_owner_absoloute_file_ssh_private_key_errors
		}

		set_source_current_owner_absoloute_file_ssh_public_key_errors := source_file_ssh_public_key.SetOwner(source_user, *source_user_group)
		if set_source_current_owner_absoloute_file_ssh_public_key_errors != nil {
			fmt.Println("set_source_current_owner_absoloute_file_ssh_public_key_errors")
			return set_source_current_owner_absoloute_file_ssh_public_key_errors
		}

		destination_file_authorised_keys, destination_file_authorised_keys_errors := newAbsoluteFile(destination_ssh_directory, "authorized_keys")
		if destination_file_authorised_keys_errors != nil {
			fmt.Println("destination_file_authorised_keys_errors")
			return destination_file_authorised_keys_errors
		}

		destination_file_authorised_keys_remove_errors := destination_file_authorised_keys.RemoveIfExists()
		if destination_file_authorised_keys_remove_errors != nil {
			fmt.Println("destination_file_authorised_keys_remove_errors")
			return destination_file_authorised_keys_remove_errors
		}

		destination_file_authorised_keys_touch_errors := destination_file_authorised_keys.Touch()
		if destination_file_authorised_keys_touch_errors != nil {
			fmt.Println("destination_file_authorised_keys_touch_errors")
			return destination_file_authorised_keys_touch_errors
		}

		destination_file_authorised_keys_append_errors := destination_file_authorised_keys.Append("blah")
		if destination_file_authorised_keys_append_errors != nil {
			fmt.Println("destination_file_authorised_keys_append_errors")
			return destination_file_authorised_keys_append_errors
		}

		set_destination_current_owner_authorized_keys_errors := destination_file_authorised_keys.SetOwner(destination_user, *destination_user_group)
		if set_destination_current_owner_authorized_keys_errors != nil {
			fmt.Println("set_destination_current_owner_authorized_keys_errors")
			return set_destination_current_owner_authorized_keys_errors
		}

		return nil
	}

	generate_ssh_key := func(destination_user HostUser) []error {
		source_user := getUser()
		ssh_directory, absolute_directory_errors := create_ssh_directory(source_user)
		if absolute_directory_errors != nil {
			fmt.Println("absolute_directory_errors")
			return absolute_directory_errors
		}

		destination_ssh_directory, destination_ssh_directory_errors := create_ssh_directory(destination_user.GetUser())
		if destination_ssh_directory_errors != nil {
			fmt.Println("absolute_directory_errors")
			return destination_ssh_directory_errors
		}

		create_key_errors := create_ssh_key(*ssh_directory, destination_user, *destination_ssh_directory)
		if create_key_errors != nil {
			return create_key_errors
		}

		return nil
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

