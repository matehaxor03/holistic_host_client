package host_client

import (
	validate "github.com/matehaxor03/holistic_validator/validate"
	common "github.com/matehaxor03/holistic_common/common"
	"path/filepath"
	"os"
)

type AbsoluteDirectory struct {
	Validate func() []error
	Create func() []error
	CreateIfDoesNotExist func() []error
	Exists func() bool
	GetPath func() []string
	GetPathAsString func() string
	SetOwnerRecursive func(host_user User, group Group) []error
	SetOwner func(host_user User, group Group) []error
	Chmod func(permissions int) []error
}

func newAbsoluteDirectory(path []string) (*AbsoluteDirectory, []error) {
	bashCommand := common.NewBashCommand()
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

	chmod := func(permissions int) []error {
		var errors []error
		chmod_error := os.Chmod(getPathAsString(), os.FileMode(permissions))
		if chmod_error != nil {
			errors = append(errors, chmod_error)
		}
		
		if len(errors) > 0 {
			return errors
		}

		return nil
	}

	exists := func() (bool) {
		exists := false
		if _, err := os.Stat(getPathAsString()); err == nil || os.IsExist(err) {
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

	createIfDoesNotExist := func() []error {
		if exists() {
			return nil
		}

		return create()
	}

	setOwnerRecursive := func(user User, group Group) []error {
		shell_command := "chown -R " + user.GetUsername() + ":" + group.GetGroupName() + " " + getPathAsString()
		_, std_error := bashCommand.ExecuteUnsafeCommandUsingFilesWithoutInputFile(shell_command)
		if std_error != nil {
			return std_error
		}
		return nil
	}

	setOwner := func(user User, group Group) []error {
		shell_command := "chown " + user.GetUsername() + ":" + group.GetGroupName() + " " + getPathAsString()
		_, std_error := bashCommand.ExecuteUnsafeCommandUsingFilesWithoutInputFile(shell_command)
		if std_error != nil {
			return std_error
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
		CreateIfDoesNotExist: func() []error {
			return createIfDoesNotExist()
		},
		Exists: func() (bool) {
			return exists()
		},
		GetPath: func() ([]string) {
			return getPath()
		},
		GetPathAsString: func() string {
			return getPathAsString()
		},
		SetOwnerRecursive: func(host_user User, group Group) []error {
			return setOwnerRecursive(host_user, group)
		},
		SetOwner: func(host_user User, group Group) []error {
			return setOwner(host_user, group)
		},
		Chmod: func(permissions int) []error {
			return chmod(permissions)
		},
	}
	setPath(path)

	errors := validate()

	if errors != nil {
		return nil, errors
	}

	return &x, nil
}

