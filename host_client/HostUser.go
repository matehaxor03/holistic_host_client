package host_client

import (
	validate "github.com/matehaxor03/holistic_validator/validate"
	"strings"
	"fmt"
)

type HostUser struct {
	Validate func() []error
}

func newHostUser(username string) (*HostUser, []error) {
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

		username_errors := verify.ValidateUsername(temp_username)
		if username_errors != nil {
			errors = append(errors, username_errors...)
		}

		if len(errors) > 0 {
			return errors
		}
		return nil
	}

	x := HostUser{
		Validate: func() []error {
			return validate()
		},
	}
	setUsername(username)

	errors := validate()

	if errors != nil {
		return nil, errors
	}

	return &x, nil
}

