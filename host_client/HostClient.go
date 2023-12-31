package host_client

import (
	"fmt"
	"os"
	validate "github.com/matehaxor03/holistic_validator/validate"	
)

type HostClient struct {
	Validate func() []error
	CreateHostUser func(username string) (*HostUser, []error)
	DeleteHostUser func(username string) ([]error)
	Ramdisk func(disk_name string,block_size uint64) (*Ramdisk, []error)
	GetEnviornmentVariable func(environment_variable_name string) (*string, []error)
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

	ramdisk := func(disk_name string, block_size uint64) (*Ramdisk, []error) {
		ramdisk, ramdisk_errors := newRamdisk(disk_name, block_size)
		if ramdisk_errors != nil {
			return nil, ramdisk_errors
		}

		return ramdisk, nil
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
			return ramdisk(disk_name, block_size)
		},
		GetEnviornmentVariable: func(environment_variable_name string) (*string, []error) {
			return get_environment_variable(environment_variable_name)
		},
	}
	//setHostClient(&x)

	errors := validate()

	if errors != nil {
		return nil, errors
	}

	return &x, nil
}

