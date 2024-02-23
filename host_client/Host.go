package host_client

import (
	validate "github.com/matehaxor03/holistic_validator/validate"
)

type Host struct {
	Validate      func() []error
	GetHostName   func() (string)
}

func newHost(host_name string) (*Host, []error) {
	verify := validate.NewValidator()
	this_host_name := host_name

	getHostName := func() (string) {
		return this_host_name
	}

	validate := func() []error {
		var errors []error
		temp_host_name := getHostName()
		if hostname_errors := verify.ValidateDomainName(temp_host_name); hostname_errors != nil {
			errors = append(errors, hostname_errors...)
		}

		if len(errors) > 0 {
			return errors
		}

		return nil
	}

	validate_errors := validate()

	if validate_errors != nil {
		return nil, validate_errors
	}

	return &Host{
		Validate: func() []error {
			return validate()
		},
		GetHostName: func() (string) {
			return getHostName()
		},
	}, nil
}
