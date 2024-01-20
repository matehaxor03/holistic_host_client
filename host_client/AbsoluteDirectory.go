package host_client

import (
	validate "github.com/matehaxor03/holistic_validator/validate"
	"path/filepath"
	"os"
)

type AbsoluteDirectory struct {
	Validate func() []error
	Create func() []error
	Exists func() bool
	GetPath func() []string
}

func newAbsoluteDirectory(path []string) (*AbsoluteDirectory, []error) {
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
		exists := true
		if _, stat_error := os.Stat(getPathAsString()); os.IsNotExist(stat_error) {
			exists = false
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

	x := AbsoluteDirectory{
		Validate: func() []error {
			return validate()
		},
		Create: func() []error {
			return create()
		},
		Exists: func() (bool) {
			return exists()
		},
		GetPath: func() ([]string) {
			return getPath()
		},
	}
	setPath(path)

	errors := validate()

	if errors != nil {
		return nil, errors
	}

	return &x, nil
}

