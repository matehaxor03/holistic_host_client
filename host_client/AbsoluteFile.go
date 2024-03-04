package host_client

import (
	validate "github.com/matehaxor03/holistic_validator/validate"
	common "github.com/matehaxor03/holistic_common/common"
	"os"
	"bufio"
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
	Touch func() []error
	ReadAllAsStringArray func() (*[]string, []error)
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

	touch := func() []error {
		shell_command := "touch " + getPathAsString()
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

	appendString := func(value string) []error {
		var errors []error
		file, file_error := os.OpenFile(getPathAsString(), os.O_APPEND|os.O_WRONLY, 0644)
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

	readAllAsStringArray := func() (*[]string, []error) {
		var errors []error
		var lines []string
		
		file, file_error := os.Open(getPathAsString())
		if file_error != nil {
			errors = append(errors, file_error)
			return nil, errors
		}
		
		defer file.Close()
		scanner := bufio.NewScanner(file)

		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}

		if scanner_error := scanner.Err(); scanner_error != nil {
			errors = append(errors, scanner_error)
			return nil, errors
		}

		return &lines, nil
	}

	x := AbsoluteFile{
		Validate: func() []error {
			return validate()
		},
		Create: func() []error {
			return create()
		},
		Append: func(value string) []error {
			return appendString(value)
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
		Touch: func() []error {
			return touch()
		},
		ReadAllAsStringArray: func() (*[]string, []error) {
			return readAllAsStringArray()
		},
	}

	setAbsoluteDirectory(directory)
	setFilename(filename)

	validate_errors := validate()

	if validate_errors != nil {
		return nil, validate_errors
	}

	return &x, nil
}

