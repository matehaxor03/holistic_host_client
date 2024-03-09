package host_client

import (
	validate "github.com/matehaxor03/holistic_validator/validate"
	common "github.com/matehaxor03/holistic_common/common"
	json "github.com/matehaxor03/holistic_json/json"	
	"strings"
	"fmt"
	"strconv"
)

type User struct {
	Validate func() []error
	Create func() []error
	Delete func() []error
	DeleteIfExists func() []error
	Exists func() (*bool, []error)
	CreateHomeDirectoryAbsoluteDirectory func(absolute_directory AbsoluteDirectory) []error
	GetHomeDirectoryAbsoluteDirectory func() (*AbsoluteDirectory, []error)
	EnableBinBash func() []error
	GetUsername func() string
	SetUniqueId func(unique_id uint64) []error
	GetUniqueId func() (*uint64, []error)
	SetPrimaryGroupId func(primary_group_id uint64) []error
	GetPrimaryGroupId func() (*uint64, []error)
	GetPrimaryGroup func() (*Group, []error)
	SetPassword func(password string) []error
}

func newUser(username string) (*User, []error) {
	bashCommand := common.NewBashCommand()
	verify := validate.NewValidator()
	var this_username string
	var this_home_directory *AbsoluteDirectory

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

	getAttribute := func(attribute string) (*json.Value,[]error) {
		var errors []error
		//todo validate attribute

		shell_command := "dscl . read /Users/" + getUsername() + " " + attribute
		std_outs, std_errors := bashCommand.ExecuteUnsafeCommandUsingFilesWithoutInputFile(shell_command)
		
		if std_errors != nil {
			std_errors = append([]error{fmt.Errorf("%s", shell_command)} , std_errors...)
			errors = append(errors, std_errors...)
		}

		if len(errors) > 0 {
			for _, err := range errors {
				if strings.Contains(fmt.Sprintf("%s", err), "RecordNotFound") {
					return nil, nil
				}
			}
			return nil, errors
		} else {
			for _, std_out := range std_outs {
				if strings.Contains(std_out, attribute + ": ") {
					raw_path := strings.TrimPrefix(std_out, attribute + ":")
					raw_path = strings.TrimSpace(raw_path)
					json_value := json.NewValue(raw_path)
					return json_value, nil
				}
			}

			errors = append(errors, fmt.Errorf("unable to determine if attribute" + attribute + " or not"))
			return nil, errors
		}
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

	getHomeDirectoryAbsoluteDirectory := func() (*AbsoluteDirectory,[]error) {
		if this_home_directory != nil {
			return this_home_directory, nil
		}

		var errors []error
		shell_command := "dscl . read /Users/" + getUsername() + " NFSHomeDirectory"
		std_outs, std_errors := bashCommand.ExecuteUnsafeCommandUsingFilesWithoutInputFile(shell_command)
		
		if std_errors != nil {
			std_errors = append([]error{fmt.Errorf("%s", shell_command)} , std_errors...)
			errors = append(errors, std_errors...)
		}

		if len(errors) > 0 {
			for _, err := range errors {
				if strings.Contains(fmt.Sprintf("%s", err), "RecordNotFound") {
					return nil, nil
				}
			}
			return nil, errors
		} else {
			for _, std_out := range std_outs {
				if strings.Contains(std_out, "NFSHomeDirectory: ") {
					raw_path := strings.TrimPrefix(std_out, "NFSHomeDirectory:")
					raw_path = strings.TrimSpace(raw_path)
					parts := strings.Split(raw_path, "/")
					absolute_directory, absolute_directory_errors := newAbsoluteDirectory(parts[1:])
					if absolute_directory_errors != nil {
						return nil, absolute_directory_errors
					}
					this_home_directory = absolute_directory
					return absolute_directory, nil
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

	deleteIfExists := func() []error {
		exists, exists_errors := exists()
		if exists_errors != nil {
			return exists_errors
		}

		if !*exists {
			return nil
		}

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

	//expect -c 'spawn ssh -vv -i ~/.ssh/holisticxyz_b3047_ holisticxyz_b3047_@127.0.0.1 "whoami"; expect "assword:"; send "passowrdhere\r"; interact'

	setPassword := func(password string) []error {
		//todo validate input
		shell_command := "dscl . -create /Users/" + getUsername() + " Password '" + password + "'"
		_, std_errors := bashCommand.ExecuteUnsafeCommandUsingFilesWithoutInputFile(shell_command)
		if std_errors != nil {
			std_errors = append([]error{fmt.Errorf("%s", shell_command)} , std_errors...)
			return std_errors
		}
		return nil
	}

	getUniqueId := func() (*uint64, []error) {
		attribute_value, getAttribute_errors := getAttribute("UniqueId")
		if getAttribute_errors != nil {
			return nil, getAttribute_errors
		}
		
		uint64_value, uint64_value_errors := attribute_value.GetUInt64()
		if uint64_value_errors != nil {
			return nil, uint64_value_errors
		}

		return uint64_value, nil
	}

	getPrimaryGroupId := func() (*uint64, []error) {
		attribute_value, getAttribute_errors := getAttribute("PrimaryGroupID")
		if getAttribute_errors != nil {
			return nil, getAttribute_errors
		}
		
		uint64_value, uint64_value_errors := attribute_value.GetUInt64()
		if uint64_value_errors != nil {
			return nil, uint64_value_errors
		}

		return uint64_value, nil
	}

	getPrimaryGroup := func() (*Group, []error) {
		primary_group_id, primary_group_id_errors := getPrimaryGroupId()
		if primary_group_id_errors != nil {
			return nil, primary_group_id_errors
		}
		
		
		var errors []error
		shell_command := "dscl . -search /Groups PrimaryGroupID " +  strconv.FormatUint(*primary_group_id, 10)
		std_outs, std_errors := bashCommand.ExecuteUnsafeCommandUsingFilesWithoutInputFile(shell_command)
		
		if std_errors != nil {
			std_errors = append([]error{fmt.Errorf("%s", shell_command)} , std_errors...)
			errors = append(errors, std_errors...)
		}

		if len(errors) > 0 {
			return nil, errors
		} else {
			for _, std_out := range std_outs {
				if strings.Contains(std_out, "PrimaryGroupID") {
					index := strings.Index(std_out, "PrimaryGroupID")
					group_name := strings.TrimSpace(std_out[:index-1])
					group, group_errors := newGroup(group_name)
					if group_errors != nil {
						return nil, group_errors
					}
					return group, nil
				}
			}

			errors = append(errors, fmt.Errorf("unable to determine if group exists or not"))
			return nil, errors
		}
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
		//todo clone absolute_directory
		
		shell_command := "dscl . -create /Users/" + getUsername() + " NFSHomeDirectory " + absolute_directory.GetPathAsString()
		_, std_errors := bashCommand.ExecuteUnsafeCommandUsingFilesWithoutInputFile(shell_command)
		if std_errors != nil {
			std_errors = append([]error{fmt.Errorf("%s", shell_command)} , std_errors...)
			return std_errors
		}
		this_home_directory = &absolute_directory
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

	x := User{
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
		DeleteIfExists: func() []error {
			return deleteIfExists()
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
		SetPassword: func(password string) []error {
			return setPassword(password)
		},
		GetUniqueId: func() (*uint64, []error) {
			return getUniqueId()
		},
		SetPrimaryGroupId: func(primary_group_id uint64) []error {
			return setPrimaryGroupId(primary_group_id)
		},
		GetPrimaryGroupId: func() (*uint64, []error) {
			return getPrimaryGroupId()
		},
		GetPrimaryGroup: func() (*Group, []error) {
			return getPrimaryGroup()
		},
		GetHomeDirectoryAbsoluteDirectory: func() (*AbsoluteDirectory, []error) {
			return getHomeDirectoryAbsoluteDirectory()
		},
	}
	setUsername(username)

	errors := validate()

	if errors != nil {
		return nil, errors
	}

	return &x, nil
}

