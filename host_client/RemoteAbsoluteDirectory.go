package host_client

import (
	validate "github.com/matehaxor03/holistic_validator/validate"
	"fmt"
	"path/filepath"
)

type RemoteAbsoluteDirectory struct {
	Validate func() []error
	Create func() []error
	CreateIfDoesNotExist func() []error
	DeleteIfExists func() []error
	Exists func() (*bool, []error)
	GetPath func() []string
	GetPathAsString func() string
}

func newRemoteAbsoluteDirectory(source_user User, host_user HostUser, path []string) (*RemoteAbsoluteDirectory, []error) {
	verify := validate.NewValidator()
	var this_path []string
	var this_source_user User
	var this_host_user HostUser

	setPath := func(path []string) {
		this_path = path
	}

	getPath := func() []string {
		return this_path
	}

	setHostUser := func(host_user HostUser) {
		this_host_user = host_user
	}

	getHostUser := func() HostUser {
		return this_host_user
	}

	setSourceUser := func(source_user User) {
		this_source_user = source_user
	}

	getSourceUser := func() User {
		return this_source_user
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

	exists := func() (*bool, []error) {
		var errors []error
		temp_source_user := getSourceUser()
		temp_host_user := getHostUser()
		shell_command := "ssh -i ~/.ssh/" + temp_host_user.GetFullyQualifiedUsername() + " " + host_user.GetFullyQualifiedUsername() + " '[ -d " + getPathAsString() + " ] && echo true || echo false'"
		std_outs, std_errors := temp_source_user.ExecuteUnsafeCommandUsingFilesWithoutInputFile(shell_command)
		
		if std_errors != nil {
			std_errors = append([]error{fmt.Errorf("%s", shell_command)} , std_errors...)
			return nil, std_errors
		}

		if len(std_outs) != 0 {
			errors = append(errors, fmt.Errorf("RemoteAboslouteDirectory exists std_out has zero output"))
			return nil, errors
		}

		var result bool
		if std_outs[0] == "true" {
			result = true
			return &result, nil
		}

		if std_outs[0] == "false" {
			result = false
			return &result, nil
		}

		errors = append(errors, fmt.Errorf("RemoteAboslouteDirectory exists unable to determine if it exists"))
		return nil, errors
	}

	create := func() []error {
		temp_source_user := getSourceUser()
		temp_host_user := getHostUser()
		shell_command := "ssh -i ~/.ssh/" + temp_host_user.GetFullyQualifiedUsername() + " " + host_user.GetFullyQualifiedUsername() + "'mkdir " + getPathAsString() + "'"
		_, std_errors := temp_source_user.ExecuteUnsafeCommandUsingFilesWithoutInputFile(shell_command)
		
		if std_errors != nil {
			return std_errors
		}

		return nil
	}

	delete := func() []error {
		temp_source_user := getSourceUser()
		temp_host_user := getHostUser()
		shell_command := "ssh -i ~/.ssh/" + temp_host_user.GetFullyQualifiedUsername() + " " + host_user.GetFullyQualifiedUsername() + "'rm -fr " + getPathAsString() + "'"
		_, std_errors := temp_source_user.ExecuteUnsafeCommandUsingFilesWithoutInputFile(shell_command)
		
		if std_errors != nil {
			return std_errors
		}

		return nil
	}

	createIfDoesNotExist := func() []error {
		exists, exists_errors := exists()
		if exists_errors != nil {
			return exists_errors
		}

		if *exists {
			return nil
		}

		return create()
	}

	deleteIfExists := func() []error {
		exists, exists_errors := exists()
		if exists_errors != nil {
			return exists_errors
		}

		if !*exists {
			return nil
		}

		return delete()
	}

	x := RemoteAbsoluteDirectory{
		Validate: func() []error {
			return validate()
		},
		Create: func() []error {
			return create()
		},
		CreateIfDoesNotExist: func() []error {
			return createIfDoesNotExist()
		},
		DeleteIfExists: func() []error {
			return deleteIfExists()
		},
		Exists: func() (*bool, []error) {
			return exists()
		},
		GetPath: func() ([]string) {
			return getPath()
		},
		GetPathAsString: func() string {
			return getPathAsString()
		},
	}
	setPath(path)
	setHostUser(host_user)
	setSourceUser(source_user)

	errors := validate()

	if errors != nil {
		return nil, errors
	}

	return &x, nil
}

