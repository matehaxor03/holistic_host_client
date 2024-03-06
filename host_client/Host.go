package host_client

import (
	validate "github.com/matehaxor03/holistic_validator/validate"
	common "github.com/matehaxor03/holistic_common/common"
	"strings"
)

type Host struct {
	Validate      func() []error
	GetHostName   func() (string)
	EnableSSH 	  func() []error
	GetLocalhostFingerprints func() (*[]string, []error)
}

func newHost(host_name string) (*Host, []error) {
	bashCommand := common.NewBashCommand()
	verify := validate.NewValidator()
	this_host_name := host_name

	getHostName := func() (string) {
		return this_host_name
	}

	enableSSH := func() []error {
		shell_command := "systemsetup -setremotelogin on"
		_, std_error := bashCommand.ExecuteUnsafeCommandUsingFilesWithoutInputFile(shell_command)
		if std_error != nil {
			return std_error
		}
		return nil
	}

	getLocalhostFingerprints := func() (*[]string, []error) {
		var fingerprints []string
		
		{
			shell_command := "ssh-keyscan 127.0.0.1"
			std_out, std_error := bashCommand.ExecuteUnsafeCommandUsingFilesWithoutInputFile(shell_command)
			if std_error != nil {
				return nil, std_error
			} 
			
			for _, s := range std_out {
				temp := strings.TrimSpace(s)
				if strings.HasPrefix(temp, "#") {
					continue
				}
				fingerprints = append(fingerprints, temp)
			}
		}

		{
			shell_command := "ssh-keyscan localhost"
			std_out, std_error := bashCommand.ExecuteUnsafeCommandUsingFilesWithoutInputFile(shell_command)
			if std_error != nil {
				return nil, std_error
			} 
			
			for _, s := range std_out {
				temp := strings.TrimSpace(s)
				if strings.HasPrefix(temp, "#") {
					continue
				}
				fingerprints = append(fingerprints, temp)
			}
		}

		return &fingerprints, nil
	}

	validate := func() []error {
		var errors []error
		temp_host_name := getHostName()
		if hostname_errors := verify.ValidateDomainName(temp_host_name); hostname_errors != nil {
			errors = append(errors, hostname_errors...)
		}

		if len(errors) > 0 {
			return errors
		}

		return nil
	}

	validate_errors := validate()

	if validate_errors != nil {
		return nil, validate_errors
	}

	return &Host{
		Validate: func() []error {
			return validate()
		},
		GetHostName: func() (string) {
			return getHostName()
		},
		EnableSSH: func() []error {
			return enableSSH()
		},
		GetLocalhostFingerprints: func() (*[]string, []error) {
			return getLocalhostFingerprints()
		},
	}, nil
}
