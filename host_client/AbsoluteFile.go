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
	RemoveIfExists func() []error
	GetPath func() []string
	GetPathAsString func() string
	GetAbsoluteDirectory func() AbsoluteDirectory
	GetFilename func() string
	SetOwner func(host_user User, group Group) []error
	Append func(value string) []error
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
		return getAbsoluteDirectory().GetPathAsString() + "/" + getFilename()
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
		var errors []error
		create_file, create_file_errors := os.Create(getPathAsString())
		if create_file_errors != nil {
			errors = append(errors, create_file_errors)
			return errors
		}
		defer create_file.Close()
		return nil
	}

	remove := func() []error {
		var errors []error
		remove_error := os.Remove(getPathAsString())
		if remove_error != nil {
			errors = append(errors, remove_error)
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

	removeIfExists := func() []error {
		if !exists() {
			return nil
		}

		return remove()
	}

	append := func(value string) []error {
		var errors []error
		file, file_error := os.OpenFile(getPathAsString(), os.O_APPEND, 0644)
		if file_error != nil {
			errors = append(errors, file_error)
			return errors
		}

		defer file.Close()
		if _, append_error := file.WriteString(value); append_error != nil {
			errors = append(errors, append_error)
			return errors
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
		Append: func(value string) []error {
			return append(value)
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
		RemoveIfExists: func() []error {
			return removeIfExists()
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

