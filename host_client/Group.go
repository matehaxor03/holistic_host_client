package host_client

import (
	validate "github.com/matehaxor03/holistic_validator/validate"
	common "github.com/matehaxor03/holistic_common/common"
	"strings"
	"fmt"
	"strconv"
)

type Group struct {
	Validate func() []error
	Create func() []error
	Delete func() []error
	GetGroupName func() string
	Exists func() (*bool, []error)
	SetUniqueId func(unique_id uint64) []error
	AddUser func(user HostUser) []error
}

func newGroup(group_name string) (*Group, []error) {
	bashCommand := common.NewBashCommand()
	verify := validate.NewValidator()
	var this_group_name string

	setGroupName := func(group_name string) {
		this_group_name = group_name
	}

	getGroupName := func() string {
		return this_group_name
	}

	validate := func() []error {
		var errors []error
		temp_group_name := getGroupName()

		if !strings.HasPrefix(temp_group_name, "holisticxyz_") {
			errors = append(errors, fmt.Errorf("temp_group_name does not start with holisticxyz_"))
		}

		if !strings.HasSuffix(temp_group_name, "_") {
			errors = append(errors, fmt.Errorf("temp_group_name does not end with _"))
		}

		group_name_errors := verify.ValidateUsername(temp_group_name)
		if group_name_errors != nil {
			errors = append(errors, group_name_errors...)
		}

		if len(errors) > 0 {
			return errors
		}
		return nil
	}

	exists := func() (*bool, []error) {
		var errors []error
		shell_command := "dscl . read /Groups/" + getGroupName() + " RecordName"
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

			errors = append(errors, fmt.Errorf("unable to determine if the group exists or not"))
			return nil, errors
		}
	}

	create := func() []error {
		shell_command := "dscl . -create /Groups/" + getGroupName()
		_, std_error := bashCommand.ExecuteUnsafeCommandUsingFilesWithoutInputFile(shell_command)
		if std_error != nil {
			return std_error
		}
		return nil
	}

	delete := func() []error {
		shell_command := "dscl . -delete /Groups/" + getGroupName()
		_, std_error := bashCommand.ExecuteUnsafeCommandUsingFilesWithoutInputFile(shell_command)
		if std_error != nil {
			return std_error
		}
		return nil
	}

	setUniqueId := func(unique_id uint64) []error {
		shell_command := "dscl . -create /Groups/" + getGroupName() + " gid " + strconv.FormatUint(unique_id, 10)
		_, std_error := bashCommand.ExecuteUnsafeCommandUsingFilesWithoutInputFile(shell_command)
		if std_error != nil {
			return std_error
		}
		return nil
	}

	addUser := func(user HostUser) []error {
		shell_command := "dscl . append /Groups/" + getGroupName() + " GroupMembership " + user.GetUsername()
		_, std_error := bashCommand.ExecuteUnsafeCommandUsingFilesWithoutInputFile(shell_command)
		if std_error != nil {
			return std_error
		}
		return nil
	}

	x := Group{
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
		GetGroupName: func() string {
			return getGroupName()
		},
		SetUniqueId: func(unique_id uint64) []error {
			return setUniqueId(unique_id)
		},
		AddUser: func(user HostUser) []error {
			return addUser(user)
		},
	}
	setGroupName(group_name)

	errors := validate()

	if errors != nil {
		return nil, errors
	}

	return &x, nil
}

