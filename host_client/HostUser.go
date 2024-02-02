package host_client

import (
	validate "github.com/matehaxor03/holistic_validator/validate"
	common "github.com/matehaxor03/holistic_common/common"
	"strings"
	"fmt"
)

type HostUser struct {
	Validate func() []error
	Create func() []error
	Delete func() []error
	Exists func() (*bool, []error)
	CreateHomeDirectoryAbsoluteDirectory func(absolute_directory AbsoluteDirectory) []error
	EnableBinBash func() []error
	GetUsername func() string
}

func newHostUser(username string) (*HostUser, []error) {
	bashCommand := common.NewBashCommand()
	verify := validate.NewValidator()
	var this_username string

	setUsername := func(username string) {
		this_username = username
	}

	getUsername := func() string {
		return this_username
	}

	validate := func() []error {
		var errors []error
		temp_username := getUsername()

		if !strings.HasPrefix(temp_username, "holisticxyz") {
			errors = append(errors, fmt.Errorf("username does not start with holisticxyz"))
		}

		username_errors := verify.ValidateUsername(temp_username)
		if username_errors != nil {
			errors = append(errors, username_errors...)
		}

		if len(errors) > 0 {
			return errors
		}
		return nil
	}

	exists := func() (*bool, []error) {
		var errors []error
		shell_command := "dscl . read /Users/" + getUsername() + " RecordName"
		std_outs, std_errors := bashCommand.ExecuteUnsafeCommandUsingFilesWithoutInputFile(shell_command)
		
		result := false
		if std_errors != nil {
			errors = append(errors, std_errors...)
		}

		if len(errors) > 0 {
			for _, err := range errors {
				if strings.Contains(fmt.Sprintf("%s", err), "RecordNotFound") {
					result = false
					return &result, nil
				}
			}

			return nil, errors
		} else {
			for _, std_out := range std_outs {
				if strings.Contains(std_out, "RecordName:") {
					result = true
					return &result, nil
				}
			}

			errors = append(errors, fmt.Errorf("unable to determine if the user exists or not"))
			return nil, errors
		}
	}

	create := func() []error {
		shell_command := "dscl . -create /Users/" + getUsername()
		_, std_error := bashCommand.ExecuteUnsafeCommandUsingFilesWithoutInputFile(shell_command)
		if std_error != nil {
			return std_error
		}
		return nil
	}

	delete := func() []error {
		shell_command := "dscl . -delete /Users/" + getUsername()
		_, std_error := bashCommand.ExecuteUnsafeCommandUsingFilesWithoutInputFile(shell_command)
		if std_error != nil {
			return std_error
		}
		return nil
	}

	createHomeDirectoryAbsoluteDirectory := func(absolute_directory AbsoluteDirectory) []error {
		shell_command := "dscl . -create /Users/" + getUsername() + " NFSHomeDirectory " + absolute_directory.GetPathAsString()
		_, std_error := bashCommand.ExecuteUnsafeCommandUsingFilesWithoutInputFile(shell_command)
		if std_error != nil {
			return std_error
		}
		return nil
	}

	enableBinBash := func() []error {
		var errors []error
		exists, exists_error := exists()
		if exists_error != nil {
			return exists_error
		}

		if !*exists {
			errors = append(errors, fmt.Errorf("user does not exist"))
			return errors
		}

		shell_command := "dscl . -create /Users/" + getUsername() + " UserShell /bin/bash"
		_, std_error := bashCommand.ExecuteUnsafeCommandUsingFilesWithoutInputFile(shell_command)
		if std_error != nil {
			return std_error
		}
		return nil
	}

	x := HostUser{
		Validate: func() []error {
			return validate()
		},
		Create: func() []error {
			return create()
		},
		Exists: func() (*bool, []error) {
			return exists()
		},
		Delete: func() []error {
			return delete()
		},
		CreateHomeDirectoryAbsoluteDirectory: func(absolute_directory AbsoluteDirectory) []error {
			return createHomeDirectoryAbsoluteDirectory(absolute_directory)
		},
		EnableBinBash: func() []error {
			return enableBinBash()
		},
		GetUsername: func() string {
			return getUsername()
		},
	}
	setUsername(username)

	errors := validate()

	if errors != nil {
		return nil, errors
	}

	return &x, nil
}

