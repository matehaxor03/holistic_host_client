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
	CreateHostUser func(username string) (*HostUser, []error)
	DeleteHostUser func(username string) ([]error)
	Ramdisk func(disk_name string,block_size uint64) (*Ramdisk, []error)
	GetEnviornmentVariable func(environment_variable_name string) (*string, []error)
	GetEnviornmentVariableValue func(environment_variable_name string) (*json.Value, []error)
	AbsoluteDirectory func(path []string) (*AbsoluteDirectory, []error)
}

func NewHostClient() (*HostClient, []error) {
	verify := validate.NewValidator()
	
	createHostUser := func(username string) (*HostUser, []error) {
		host_user, host_user_errors := newHostUser(username)
		
		if host_user_errors != nil {
			return nil, host_user_errors
		}

		create_errors := host_user.Create()
		
		if create_errors != nil {
			return nil, create_errors
		}

		return host_user, nil
	}

	deleteHostUser := func(username string) ([]error) {
		host_user, host_user_errors := newHostUser(username)
		
		if host_user_errors != nil {
			return host_user_errors
		}

		delete_errors := host_user.Delete()
		
		if delete_errors != nil {
			return delete_errors
		}

		return nil
	}

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
		CreateHostUser: func(username string) (*HostUser, []error) {
			return createHostUser(username)
		},
		DeleteHostUser: func(username string) ([]error) {
			return deleteHostUser(username)
		},
		Ramdisk: func(disk_name string, block_size uint64) (*Ramdisk, []error) {
			return newRamdisk(disk_name, block_size)
		},
		HostUser: func(username string) (*HostUser, []error) {
			return newHostUser(username)
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

