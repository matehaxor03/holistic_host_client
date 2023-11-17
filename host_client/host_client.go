package host_client

import (
	validate "github.com/matehaxor03/holistic_validator/validate"
)

type HostClient struct {
	Validate func() []error
}

func NewHostClient() (*HostClient, []error) {
	verify := validate.NewValidator()
	//var this_host_client *HostClient

	/*setHostClient := func(host_client *HostClient) {
		this_host_client = host_client
	}*/

	/*
	getHostClient := func() *HostClient {
		return this_host_client
	}*/

	validate := func() []error {
		return verify.ValidateUsername("asdfasf")
	}

	x := HostClient{
		Validate: func() []error {
			return validate()
		},
	}
	//setHostClient(&x)

	errors := validate()

	if errors != nil {
		return nil, errors
	}

	return &x, nil
}

