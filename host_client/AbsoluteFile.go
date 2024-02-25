package host_client

import (
	validate "github.com/matehaxor03/holistic_validator/validate"
	common "github.com/matehaxor03/holistic_common/common"
	"os"
)

type AbsoluteFile struct {
	Validate func() []error
	Create func() []error
	Exists func() bool
	GetPath func() []string
	GetPathAsString func() string
	GetAbsoluteDirectory func() AbsoluteDirectory
	GetFilename func() string
	SetOwner func(host_user User, group Group) []error
}

func newAbsoluteFile(directory AbsoluteDirectory, filename string) (*AbsoluteFile, []error) {
	bashCommand := common.NewBashCommand()
	verify := validate.NewValidator()
	var this_absolute_directory AbsoluteDirectory
	var this_filename string

	setAbsoluteDirectory := func(dir AbsoluteDirectory) {
		this_absolute_directory = dir
	}

	getAbsoluteDirectory := func() AbsoluteDirectory {
		return this_absolute_directory
	}

	getFilename := func() string {
		return this_filename
	}

	setFilename := func(in_filename string) {
		this_filename = in_filename
	}

	getPathAsString := func() string {
		return 	"/" + getAbsoluteDirectory().GetPathAsString() + "/" + getFilename()
	}

	validate := func() []error {
		filename_errors :=  verify.ValidateFileName(getFilename())
		if filename_errors != nil {
			return filename_errors
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
		//todo fix
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

	setOwner := func(user User, group Group) []error {
		shell_command := "chown " + user.GetUsername() + ":" + group.GetGroupName() + " " + getPathAsString()
		_, std_error := bashCommand.ExecuteUnsafeCommandUsingFilesWithoutInputFile(shell_command)
		if std_error != nil {
			return std_error
		}
		return nil
	}

	x := AbsoluteFile{
		Validate: func() []error {
			return validate()
		},
		Create: func() []error {
			return create()
		},
		Exists: func() (bool) {
			return exists()
		},
		GetPathAsString: func() string {
			return getPathAsString()
		},
		GetFilename: func() string {
			return getFilename()
		},
		GetAbsoluteDirectory: func() AbsoluteDirectory {
			return getAbsoluteDirectory()
		},
		SetOwner: func(host_user User, group Group) []error {
			return setOwner(host_user, group)
		},
	}
	setAbsoluteDirectory(directory)
	setFilename(filename)

	errors := validate()

	if errors != nil {
		return nil, errors
	}

	return &x, nil
}

