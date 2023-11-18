package host_client

import (
)

type HostClient struct {
	Validate func() []error
	CreateHostUser func(username string) (*HostUser, []error)
}

func NewHostClient() (*HostClient, []error) {
	
	createHostUser := func(username string) (*HostUser, []error) {
		return newHostUser(username)
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
	}
	//setHostClient(&x)

	errors := validate()

	if errors != nil {
		return nil, errors
	}

	return &x, nil
}

