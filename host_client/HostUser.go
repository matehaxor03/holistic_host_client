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

		if !strings.HasPrefix(temp_username, "holisticxyz_") {
			errors = append(errors, fmt.Errorf("username does not start with holisticxyz_"))
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

	x := HostUser{
		Validate: func() []error {
			return validate()
		},
		Create: func() []error {
			return create()
		},
		Delete: func() []error {
			return delete()
		},
	}
	setUsername(username)

	errors := validate()

	if errors != nil {
		return nil, errors
	}

	return &x, nil
}

