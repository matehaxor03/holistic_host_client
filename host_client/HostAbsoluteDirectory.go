package host_client

import (
	validate "github.com/matehaxor03/holistic_validator/validate"
	"path/filepath"
	"os"
)

type HostAbsoluteDirectory struct {
	Validate func() []error
	Create func() []error
	Exists func() bool
}

func newHostAbsoluteDirectory(path []string) (*HostAbsoluteDirectory, []error) {
	verify := validate.NewValidator()
	var this_path []string

	setPath := func(path []string) {
		this_path = path
	}

	getPath := func() []string {
		return this_path
	}

	getPathAsString := func() string {
		return 	"/" + filepath.Join(getPath()...)
	}

	validate := func() []error {
		var errors []error
		temp_path := getPath()

		for _, s := range temp_path {
			directory_name_errors := verify.ValidateDirectoryName(s)
			if directory_name_errors != nil {
				errors = append(errors, directory_name_errors...)
			}
		}

		if len(errors) > 0 {
			return errors
		}
		return nil
	}

	exists := func() (bool) {
		exists := false
		if _, err := os.Stat(getPathAsString()); !os.IsNotExist(err) {
			exists = true
		}
		return exists
	}

	create := func() []error {
		var errors []error
		permissions := int(0700)
		create_directory_error := os.MkdirAll(getPathAsString(), os.FileMode(permissions))
		if create_directory_error != nil {
			errors = append(errors, create_directory_error)
		}
		
		if len(errors) > 0 {
			return errors
		}

		return nil
	}

	x := HostAbsoluteDirectory{
		Validate: func() []error {
			return validate()
		},
		Create: func() []error {
			return create()
		},
		Exists: func() (bool) {
			return exists()
		},
	}
	setPath(path)

	errors := validate()

	if errors != nil {
		return nil, errors
	}

	return &x, nil
}

