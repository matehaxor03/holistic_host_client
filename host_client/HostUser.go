package host_client

import (
	validate "github.com/matehaxor03/holistic_validator/validate"
	common "github.com/matehaxor03/holistic_common/common"
	"strings"
	"fmt"
	"strconv"
)

type HostUser struct {
	Validate func() []error
	Create func() []error
	Delete func() []error
	Exists func() (*bool, []error)
	CreateHomeDirectoryAbsoluteDirectory func(absolute_directory AbsoluteDirectory) []error
	EnableBinBash func() []error
	GetUsername func() string
	SetUniqueId func(unique_id uint64) []error
	SetPrimaryGroupId func(primary_group_id uint64) []error
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

		if !strings.HasSuffix(temp_username, "_") {
			errors = append(errors, fmt.Errorf("username does not end with _"))
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
			std_errors = append([]error{fmt.Errorf("%s", shell_command)} , std_errors...)
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
		_, std_errors := bashCommand.ExecuteUnsafeCommandUsingFilesWithoutInputFile(shell_command)
		if std_errors != nil {
			std_errors = append([]error{fmt.Errorf("%s", shell_command)} , std_errors...)
			return std_errors
		}
		return nil
	}

	delete := func() []error {
		shell_command := "dscl . -delete /Users/" + getUsername()
		_, std_errors := bashCommand.ExecuteUnsafeCommandUsingFilesWithoutInputFile(shell_command)
		if std_errors != nil {
			std_errors = append([]error{fmt.Errorf("%s", shell_command)} , std_errors...)
			return std_errors
		}
		return nil
	}

	setUniqueId := func(unique_id uint64) []error {
		shell_command := "dscl . -create /Users/" + getUsername() + " UniqueID " + strconv.FormatUint(unique_id, 10)
		_, std_errors := bashCommand.ExecuteUnsafeCommandUsingFilesWithoutInputFile(shell_command)
		if std_errors != nil {
			std_errors = append([]error{fmt.Errorf("%s", shell_command)} , std_errors...)
			return std_errors
		}
		return nil
	}

	setPrimaryGroupId := func(primary_group_id uint64) []error {
		shell_command := "dscl . -create /Users/" + getUsername() + " PrimaryGroupID " + strconv.FormatUint(primary_group_id, 10)
		_, std_errors := bashCommand.ExecuteUnsafeCommandUsingFilesWithoutInputFile(shell_command)
		if std_errors != nil {
			std_errors = append([]error{fmt.Errorf("%s", shell_command)} , std_errors...)
			return std_errors
		}
		return nil
	}

	createHomeDirectoryAbsoluteDirectory := func(absolute_directory AbsoluteDirectory) []error {
		shell_command := "dscl . -create /Users/" + getUsername() + " NFSHomeDirectory " + absolute_directory.GetPathAsString()
		_, std_errors := bashCommand.ExecuteUnsafeCommandUsingFilesWithoutInputFile(shell_command)
		if std_errors != nil {
			std_errors = append([]error{fmt.Errorf("%s", shell_command)} , std_errors...)
			return std_errors
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
		_, std_errors := bashCommand.ExecuteUnsafeCommandUsingFilesWithoutInputFile(shell_command)
		if std_errors != nil {
			std_errors = append([]error{fmt.Errorf("%s", shell_command)} , std_errors...)
			return std_errors
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
		SetUniqueId: func(unique_id uint64) []error {
			return setUniqueId(unique_id)
		},
		SetPrimaryGroupId: func(primary_group_id uint64) []error {
			return setPrimaryGroupId(primary_group_id)
		},
	}
	setUsername(username)

	errors := validate()

	if errors != nil {
		return nil, errors
	}

	return &x, nil
}

