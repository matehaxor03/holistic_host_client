package host_client

import (
	"fmt"
	"os"
	validate "github.com/matehaxor03/holistic_validator/validate"	
	json "github.com/matehaxor03/holistic_json/json"	
)

type HostClient struct {
	Validate func() []error
	HostUser func(username string) (*HostUser, []error)
	Group func(group_name string) (*Group, []error)
	Ramdisk func(disk_name string,block_size uint64) (*Ramdisk, []error)
	GetEnviornmentVariable func(environment_variable_name string) (*string, []error)
	GetEnviornmentVariableValue func(environment_variable_name string) (*json.Value, []error)
	AbsoluteDirectory func(path []string) (*AbsoluteDirectory, []error)
}

func NewHostClient() (*HostClient, []error) {
	verify := validate.NewValidator()

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
		HostUser: func(username string) (*HostUser, []error) {
			return newHostUser(username)
		},
		Group: func(group_name string) (*Group, []error) {
			return newGroup(group_name)
		},
		AbsoluteDirectory: func(path []string) (*AbsoluteDirectory, []error) {
			return newAbsoluteDirectory(path)
		},
		GetEnviornmentVariable: func(environment_variable_name string) (*string, []error) {
			return get_environment_variable(environment_variable_name)
		},
		GetEnviornmentVariableValue: func(environment_variable_name string) (*json.Value, []error) {
			return get_environment_variable_value(environment_variable_name)
		},
	}
	errors := validate()

	if errors != nil {
		return nil, errors
	}

	return &x, nil
}

