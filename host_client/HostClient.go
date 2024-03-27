package host_client

import (
	"fmt"
	"os"
	validate "github.com/matehaxor03/holistic_validator/validate"	
	json "github.com/matehaxor03/holistic_json/json"
	common "github.com/matehaxor03/holistic_common/common"	
)

type HostClient struct {
	Validate func() []error
	User func(username string) (*User, []error)
	Host func(host_name string) (*Host, []error)
	Group func(group_name string) (*Group, []error)
	Ramdisk func(disk_name string,block_size uint64) (*Ramdisk, []error)
	GetEnviornmentVariable func(environment_variable_name string) (*string, []error)
	GetEnviornmentVariableValue func(environment_variable_name string) (*json.Value, []error)
	AbsoluteDirectory func(path []string) (*AbsoluteDirectory, []error)
	HostUser func(host Host, user User) HostUser
	AbsoluteFile func(directory AbsoluteDirectory, filename string) (*AbsoluteFile, []error)
	Whoami func() (*string, []error)
}

func NewHostClient() (*HostClient, []error) {
	verify := validate.NewValidator()
	bashCommand := common.NewBashCommand()

	get_environment_variable := func(environment_variable string) (*string, []error) {
		var errors []error

		environment_variable_errors := verify.ValidateEnvironmentVariableName(environment_variable) 
		if environment_variable_errors != nil {
			return nil, environment_variable_errors
		}

		_, found := os.LookupEnv(environment_variable)
		if !found {
			errors = append(errors, fmt.Errorf("environment variable: "  + environment_variable + " does not exist"))
		}

		if len(errors) > 0 {
			return nil, errors
		}

		environment_variable_value := os.Getenv(environment_variable) 

		return &environment_variable_value, nil
	}

	get_environment_variable_value := func(environment_variable string) (*json.Value, []error) {
		env_variable_string_value, env_variable_string_value_errors :=  get_environment_variable(environment_variable)
		if env_variable_string_value_errors != nil {
			return nil, env_variable_string_value_errors
		}

		return json.NewValue(env_variable_string_value), nil
	}

	whoami := func() (*string, []error) {
		var errors []error
		shell_command := "whoami"
		std_out, std_out_errors := bashCommand.ExecuteUnsafeCommandUsingFilesWithoutInputFile(shell_command)

		if std_out_errors != nil {
			return nil, std_out_errors
		} else if len(std_out) == 0 {
			errors = append(errors, fmt.Errorf("whoami didn't return any stdout"))
			return nil, errors
		}
		result := std_out[0]

		return &result, nil
	}

	validate := func() []error {
		return nil
	}

	x := HostClient{
		Validate: func() []error {
			return validate()
		},
		Ramdisk: func(disk_name string, block_size uint64) (*Ramdisk, []error) {
			return newRamdisk(disk_name, block_size)
		},
		User: func(username string) (*User, []error) {
			return newUser(username)
		},
		Host: func(host_name string) (*Host, []error) {
			return newHost(host_name)
		},
		Group: func(group_name string) (*Group, []error) {
			return newGroup(group_name)
		},
		AbsoluteDirectory: func(path []string) (*AbsoluteDirectory, []error) {
			return newAbsoluteDirectory(path)
		},
		AbsoluteFile: func(directory AbsoluteDirectory, filename string) (*AbsoluteFile, []error) {
			return newAbsoluteFile(directory, filename)
		},
		GetEnviornmentVariable: func(environment_variable_name string) (*string, []error) {
			return get_environment_variable(environment_variable_name)
		},
		GetEnviornmentVariableValue: func(environment_variable_name string) (*json.Value, []error) {
			return get_environment_variable_value(environment_variable_name)
		},
		HostUser: func(host Host, user User) HostUser {
			return newHostUser(host, user)
		},
		Whoami: func() (*string, []error) {
			return whoami()
		},
	}
	errors := validate()

	if errors != nil {
		return nil, errors
	}

	return &x, nil
}

