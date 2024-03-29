package host_client

import (
	validate "github.com/matehaxor03/holistic_validator/validate"
	common "github.com/matehaxor03/holistic_common/common"
	json "github.com/matehaxor03/holistic_json/json"	
	"strings"
	"fmt"
	"strconv"
	"time"
	"sync"
	"io/ioutil"
	"os"
	"bufio"
	"container/list"
	"os/exec"
)

type User struct {
	Validate func() []error
	Create func() []error
	Delete func() []error
	DeleteIfExists func() []error
	Exists func() (*bool, []error)
	CreateHomeDirectoryAbsoluteDirectory func(absolute_directory AbsoluteDirectory) []error
	GetHomeDirectoryAbsoluteDirectory func() (*AbsoluteDirectory, []error)
	GetDirectoryIOAbsoluteDirectory func() (*AbsoluteDirectory, []error)
	GetDirectoryDBAbsoluteDirectory func() (*AbsoluteDirectory, []error)
	GetDirectorySSHAbsoluteDirectory func() (*AbsoluteDirectory, []error)
	EnableBinBash func() []error
	EnableRemoteFullDiskAccess func() []error
	DisableRemoteFullDiskAccess func() []error
	GetUsername func() string
	SetUniqueId func(unique_id uint64) []error
	GetUniqueId func() (*uint64, []error)
	SetPrimaryGroupId func(primary_group_id uint64) []error
	GetPrimaryGroupId func() (*uint64, []error)
	GetPrimaryGroup func() (*Group, []error)
	SetPassword func(password string) []error
	ExecuteUnsafeCommandUsingFiles func(command string, command_data string) ([]string, []error)
	ExecuteRemoteUnsafeCommandUsingFilesWithoutInputFile func(destination_user HostUser, command string) ([]string, []error)
}

func newUser(username string) (*User, []error) {
	bashCommand := common.NewBashCommand()
	verify := validate.NewValidator()
	const maxCapacity = 10*1024*1024  
	delete_files := list.New()
	lock := &sync.RWMutex{}
	status_lock := &sync.RWMutex{}
	file_lock := &sync.RWMutex{}
	var wg sync.WaitGroup
	wakeup_lock := &sync.Mutex{}
	status := "running"
	var this_username string
	var this_home_directory *AbsoluteDirectory

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

		if !strings.HasSuffix(temp_username, "_") {
			errors = append(errors, fmt.Errorf("username does not end with _"))
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

	getAttribute := func(attribute string) (*json.Value,[]error) {
		var errors []error
		//todo validate attribute

		shell_command := "dscl . read /Users/" + getUsername() + " " + attribute
		std_outs, std_errors := bashCommand.ExecuteUnsafeCommandUsingFilesWithoutInputFile(shell_command)
		
		if std_errors != nil {
			std_errors = append([]error{fmt.Errorf("%s", shell_command)} , std_errors...)
			errors = append(errors, std_errors...)
		}

		if len(errors) > 0 {
			for _, err := range errors {
				if strings.Contains(fmt.Sprintf("%s", err), "RecordNotFound") {
					return nil, nil
				}
			}
			return nil, errors
		} else {
			for _, std_out := range std_outs {
				if strings.Contains(std_out, attribute + ": ") {
					raw_path := strings.TrimPrefix(std_out, attribute + ":")
					raw_path = strings.TrimSpace(raw_path)
					json_value := json.NewValue(raw_path)
					return json_value, nil
				}
			}

			errors = append(errors, fmt.Errorf("unable to determine if attribute" + attribute + " or not"))
			return nil, errors
		}
	}

	exists := func() (*bool, []error) {
		var errors []error
		shell_command := "dscl . read /Users/" + getUsername() + " RecordName"
		std_outs, std_errors := bashCommand.ExecuteUnsafeCommandUsingFilesWithoutInputFile(shell_command)
		
		result := false
		if std_errors != nil {
			std_errors = append([]error{fmt.Errorf("%s", shell_command)} , std_errors...)
			errors = append(errors, std_errors...)
		}

		if len(errors) > 0 {
			for _, err := range errors {
				if strings.Contains(fmt.Sprintf("%s", err), "RecordNotFound") {
					result = false
					return &result, nil
				}
			}

			return nil, errors
		} else {
			for _, std_out := range std_outs {
				if strings.Contains(std_out, "RecordName:") {
					result = true
					return &result, nil
				}
			}

			errors = append(errors, fmt.Errorf("unable to determine if the user exists or not"))
			return nil, errors
		}
	}

	getHomeDirectoryAbsoluteDirectory := func() (*AbsoluteDirectory,[]error) {
		if this_home_directory != nil {
			return this_home_directory, nil
		}

		var errors []error
		shell_command := "dscl . read /Users/" + getUsername() + " NFSHomeDirectory"
		std_outs, std_errors := bashCommand.ExecuteUnsafeCommandUsingFilesWithoutInputFile(shell_command)
		
		if std_errors != nil {
			std_errors = append([]error{fmt.Errorf("%s", shell_command)} , std_errors...)
			errors = append(errors, std_errors...)
		}

		if len(errors) > 0 {
			for _, err := range errors {
				if strings.Contains(fmt.Sprintf("%s", err), "RecordNotFound") {
					return nil, nil
				}
			}
			return nil, errors
		} else {
			for _, std_out := range std_outs {
				if strings.Contains(std_out, "NFSHomeDirectory: ") {
					raw_path := strings.TrimPrefix(std_out, "NFSHomeDirectory:")
					raw_path = strings.TrimSpace(raw_path)
					parts := strings.Split(raw_path, "/")
					absolute_directory, absolute_directory_errors := newAbsoluteDirectory(parts[1:])
					if absolute_directory_errors != nil {
						return nil, absolute_directory_errors
					}
					this_home_directory = absolute_directory
					return absolute_directory, nil
				}
			}

			errors = append(errors, fmt.Errorf("unable to determine if the user exists or not"))
			return nil, errors
		}
	}

	getDirectoryIOAbsoluteDirectory := func() (*AbsoluteDirectory, []error) {
		var errors []error
		hd, hd_errors := getHomeDirectoryAbsoluteDirectory()
		if hd_errors != nil {
			return nil, hd_errors
		} else if hd == nil {
			errors = append(errors, fmt.Errorf("home directory is nil"))
		}
		
		if len(errors) > 0 {
			return nil, errors
		}
		p := hd.GetPath()
		p = append(p, ".io")

		io_dir, io_dir_errors := newAbsoluteDirectory(p)
		if io_dir_errors != nil {
			return nil, io_dir_errors
		}

		return io_dir, nil
	}

	getDirectoryDBAbsoluteDirectory := func() (*AbsoluteDirectory, []error) {
		var errors []error
		hd, hd_errors := getHomeDirectoryAbsoluteDirectory()
		if hd_errors != nil {
			return nil, hd_errors
		} else if hd == nil {
			errors = append(errors, fmt.Errorf("home directory is nil"))
		}
		
		if len(errors) > 0 {
			return nil, errors
		}
		p := hd.GetPath()
		p = append(p, ".db")

		db_dir, db_dir_errors := newAbsoluteDirectory(p)
		if db_dir_errors != nil {
			return nil, db_dir_errors
		}

		return db_dir, nil
	}

	getDirectorySSHAbsoluteDirectory := func() (*AbsoluteDirectory, []error) {
		var errors []error
		hd, hd_errors := getHomeDirectoryAbsoluteDirectory()
		if hd_errors != nil {
			return nil, hd_errors
		} else if hd == nil {
			errors = append(errors, fmt.Errorf("home directory is nil"))
		}
		
		if len(errors) > 0 {
			return nil, errors
		}
		p := hd.GetPath()
		p = append(p, ".ssh")

		ssh_dir, ssh_dir_errors := newAbsoluteDirectory(p)
		if ssh_dir_errors != nil {
			return nil, ssh_dir_errors
		}

		return ssh_dir, nil
	}

	create := func() []error {
		shell_command := "dscl . -create /Users/" + getUsername()
		_, std_errors := bashCommand.ExecuteUnsafeCommandUsingFilesWithoutInputFile(shell_command)
		if std_errors != nil {
			std_errors = append([]error{fmt.Errorf("%s", shell_command)} , std_errors...)
			return std_errors
		}
		return nil
	}

	delete := func() []error {
		shell_command := "dscl . -delete /Users/" + getUsername()
		_, std_errors := bashCommand.ExecuteUnsafeCommandUsingFilesWithoutInputFile(shell_command)
		if std_errors != nil {
			std_errors = append([]error{fmt.Errorf("%s", shell_command)} , std_errors...)
			return std_errors
		}
		return nil
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

	setUniqueId := func(unique_id uint64) []error {
		shell_command := "dscl . -create /Users/" + getUsername() + " UniqueID " + strconv.FormatUint(unique_id, 10)
		_, std_errors := bashCommand.ExecuteUnsafeCommandUsingFilesWithoutInputFile(shell_command)
		if std_errors != nil {
			std_errors = append([]error{fmt.Errorf("%s", shell_command)} , std_errors...)
			return std_errors
		}
		return nil
	}

	//expect -c 'spawn ssh -vv -i ~/.ssh/holisticxyz_b3047_ holisticxyz_b3047_@127.0.0.1 "whoami"; expect "assword:"; send "passowrdhere\r"; interact'

	setPassword := func(password string) []error {
		//todo validate input
		shell_command := "dscl . -create /Users/" + getUsername() + " Password '" + password + "'"
		_, std_errors := bashCommand.ExecuteUnsafeCommandUsingFilesWithoutInputFile(shell_command)
		if std_errors != nil {
			std_errors = append([]error{fmt.Errorf("%s", shell_command)} , std_errors...)
			return std_errors
		}
		return nil
	}

	getUniqueId := func() (*uint64, []error) {
		attribute_value, getAttribute_errors := getAttribute("UniqueId")
		if getAttribute_errors != nil {
			return nil, getAttribute_errors
		}
		
		uint64_value, uint64_value_errors := attribute_value.GetUInt64()
		if uint64_value_errors != nil {
			return nil, uint64_value_errors
		}

		return uint64_value, nil
	}

	getPrimaryGroupId := func() (*uint64, []error) {
		attribute_value, getAttribute_errors := getAttribute("PrimaryGroupID")
		if getAttribute_errors != nil {
			return nil, getAttribute_errors
		}
		
		uint64_value, uint64_value_errors := attribute_value.GetUInt64()
		if uint64_value_errors != nil {
			return nil, uint64_value_errors
		}

		return uint64_value, nil
	}

	getPrimaryGroup := func() (*Group, []error) {
		primary_group_id, primary_group_id_errors := getPrimaryGroupId()
		if primary_group_id_errors != nil {
			return nil, primary_group_id_errors
		}
		
		
		var errors []error
		shell_command := "dscl . -search /Groups PrimaryGroupID " +  strconv.FormatUint(*primary_group_id, 10)
		std_outs, std_errors := bashCommand.ExecuteUnsafeCommandUsingFilesWithoutInputFile(shell_command)
		
		if std_errors != nil {
			std_errors = append([]error{fmt.Errorf("%s", shell_command)} , std_errors...)
			errors = append(errors, std_errors...)
		}

		if len(errors) > 0 {
			return nil, errors
		} else {
			for _, std_out := range std_outs {
				if strings.Contains(std_out, "PrimaryGroupID") {
					index := strings.Index(std_out, "PrimaryGroupID")
					group_name := strings.TrimSpace(std_out[:index-1])
					group, group_errors := newGroup(group_name)
					if group_errors != nil {
						return nil, group_errors
					}
					return group, nil
				}
			}

			errors = append(errors, fmt.Errorf("unable to determine if group exists or not"))
			return nil, errors
		}
	}

	setPrimaryGroupId := func(primary_group_id uint64) []error {
		shell_command := "dscl . -create /Users/" + getUsername() + " PrimaryGroupID " + strconv.FormatUint(primary_group_id, 10)
		_, std_errors := bashCommand.ExecuteUnsafeCommandUsingFilesWithoutInputFile(shell_command)
		if std_errors != nil {
			std_errors = append([]error{fmt.Errorf("%s", shell_command)} , std_errors...)
			return std_errors
		}
		return nil
	}

	enableRemoteFullDiskAccess := func() []error {
		shell_command := "dseditgroup -o edit -t user -a " + getUsername() + " com.apple.access_ssh"
		_, std_errors := bashCommand.ExecuteUnsafeCommandUsingFilesWithoutInputFile(shell_command)
		if std_errors != nil {
			std_errors = append([]error{fmt.Errorf("%s", shell_command)} , std_errors...)
			return std_errors
		}
		return nil
	}

	disableRemoteFullDiskAccess := func() []error {
		shell_command := "dseditgroup -o edit -t user -d " + getUsername() + " com.apple.access_ssh"
		_, std_errors := bashCommand.ExecuteUnsafeCommandUsingFilesWithoutInputFile(shell_command)
		if std_errors != nil {
			error_string := fmt.Sprintf("%s", std_errors)
			if strings.Contains(error_string, "Record was not found") {
				return nil
			}
			
			std_errors = append([]error{fmt.Errorf("%s", shell_command)} , std_errors...)
			return std_errors
		}
		return nil
	}

	createHomeDirectoryAbsoluteDirectory := func(absolute_directory AbsoluteDirectory) []error {
		//todo clone absolute_directory
		
		shell_command := "dscl . -create /Users/" + getUsername() + " NFSHomeDirectory " + absolute_directory.GetPathAsString()
		_, std_errors := bashCommand.ExecuteUnsafeCommandUsingFilesWithoutInputFile(shell_command)
		if std_errors != nil {
			std_errors = append([]error{fmt.Errorf("%s", shell_command)} , std_errors...)
			return std_errors
		}
		this_home_directory = &absolute_directory
		return nil
	}

	enableBinBash := func() []error {
		var errors []error
		exists, exists_error := exists()
		if exists_error != nil {
			return exists_error
		}

		if !*exists {
			errors = append(errors, fmt.Errorf("user does not exist"))
			return errors
		}

		shell_command := "dscl . -create /Users/" + getUsername() + " UserShell /bin/bash"
		_, std_errors := bashCommand.ExecuteUnsafeCommandUsingFilesWithoutInputFile(shell_command)
		if std_errors != nil {
			std_errors = append([]error{fmt.Errorf("%s", shell_command)} , std_errors...)
			return std_errors
		}
		return nil
	}

	get_or_set_status := func(s string) string {
		status_lock.Lock()
		defer status_lock.Unlock()
		if s == "" {
			return status
		} else {
			status = s
			return ""
		}
	}

	execute_unsafe_command_simple := func(command string) ([]error) {
		var errors []error
		cmd := exec.Command("bash", "-c", command)
		
		start_error := cmd.Start()
		if start_error != nil {
			errors = append(errors, start_error)
		}
		
		wait_error := cmd.Wait()
		if wait_error != nil {
			errors = append(errors, wait_error)
		}

		if len(errors) > 0 {
			return errors
		}

		return nil
	}

	wakeup_delete_file_processor := func() {
		wakeup_lock.Lock()
		defer wakeup_lock.Unlock()
		if get_or_set_status("") == "paused" {
			get_or_set_status("try again") 
			wg.Done()
		} else {
			get_or_set_status("try again") 
		}
	}


	get_or_set_files := func(absolute_path_filename *string, mode string) (*string, error) {
		file_lock.Lock()
		defer file_lock.Unlock()
		if mode == "push" {
			if absolute_path_filename == nil {
				return nil, fmt.Errorf("absolute_path_filename is nil")
			}
			delete_files.PushFront(absolute_path_filename)
			wakeup_delete_file_processor()
			return nil, nil
		} else if mode == "pull" {
			message := delete_files.Front()
			if message == nil {
				return nil, nil
			}
			delete_files.Remove(message)
			return message.Value.(*string), nil
		} else {
			return nil, fmt.Errorf("mode is not supported %s", mode)
		}
	}


	cleanup_files := func(input_file string, stdout_file string, std_err_file string) {
		if input_file != "" {
			get_or_set_files(&input_file, "push")
		}
		get_or_set_files(&stdout_file, "push")
		get_or_set_files(&std_err_file, "push")
	}


	process_clean_up := func() {
		for {
			get_or_set_status("running")
			time.Sleep(1 * time.Nanosecond) 
			absolute_path_filename_to_delete, absolute_path_filename_to_delete_error := get_or_set_files(nil, "pull")
			if absolute_path_filename_to_delete_error != nil {
				fmt.Println(absolute_path_filename_to_delete_error)
			} else if absolute_path_filename_to_delete != nil {
				os.Remove(*absolute_path_filename_to_delete)
			} else if get_or_set_status("") == "running" {
				wg.Add(1)
				get_or_set_status("paused")
				wg.Wait()
				get_or_set_status("running")
			}
		}
	}


	x := User{
		Validate: func() []error {
			return validate()
		},
		Create: func() []error {
			return create()
		},
		Exists: func() (*bool, []error) {
			return exists()
		},
		Delete: func() []error {
			return delete()
		},
		DeleteIfExists: func() []error {
			return deleteIfExists()
		},
		CreateHomeDirectoryAbsoluteDirectory: func(absolute_directory AbsoluteDirectory) []error {
			return createHomeDirectoryAbsoluteDirectory(absolute_directory)
		},
		EnableBinBash: func() []error {
			return enableBinBash()
		},
		GetUsername: func() string {
			return getUsername()
		},
		SetUniqueId: func(unique_id uint64) []error {
			return setUniqueId(unique_id)
		},
		SetPassword: func(password string) []error {
			return setPassword(password)
		},
		GetUniqueId: func() (*uint64, []error) {
			return getUniqueId()
		},
		SetPrimaryGroupId: func(primary_group_id uint64) []error {
			return setPrimaryGroupId(primary_group_id)
		},
		GetPrimaryGroupId: func() (*uint64, []error) {
			return getPrimaryGroupId()
		},
		GetPrimaryGroup: func() (*Group, []error) {
			return getPrimaryGroup()
		},
		GetHomeDirectoryAbsoluteDirectory: func() (*AbsoluteDirectory, []error) {
			return getHomeDirectoryAbsoluteDirectory()
		},
		EnableRemoteFullDiskAccess: func() []error {
			return enableRemoteFullDiskAccess()
		},
		DisableRemoteFullDiskAccess: func() []error {
			return disableRemoteFullDiskAccess()
		},
		GetDirectoryIOAbsoluteDirectory: func() (*AbsoluteDirectory, []error) {
			return getDirectoryIOAbsoluteDirectory()
		},
		GetDirectoryDBAbsoluteDirectory: func() (*AbsoluteDirectory, []error) {
			return getDirectoryDBAbsoluteDirectory()
		},
		GetDirectorySSHAbsoluteDirectory: func() (*AbsoluteDirectory, []error) {
			return getDirectorySSHAbsoluteDirectory()
		},
		ExecuteRemoteUnsafeCommandUsingFilesWithoutInputFile: func(destination_user HostUser, command string) ([]string, []error) {
			lock.Lock()
			defer lock.Unlock()
			var errors []error
			var stdout_lines []string

			io_absolute_directory, io_absolute_directory_errors := getDirectoryIOAbsoluteDirectory()
			if io_absolute_directory_errors != nil {
				return stdout_lines, io_absolute_directory_errors
			}

			directory := io_absolute_directory.GetPathAsString()

			uuid, _ := ioutil.ReadFile("/proc/sys/kernel/random/uuid")
			time_now := time.Now().UnixNano()
			filename_stdout := directory + "/" + fmt.Sprintf("%v%s-stdout.sql", time_now, string(uuid))
			filename_stderr := directory + "/" + fmt.Sprintf("%v%s-stderr.sql", time_now, string(uuid))
			defer cleanup_files("", filename_stdout, filename_stderr)

			ssh_directory, ssh_directory_errors := getDirectorySSHAbsoluteDirectory()
			if ssh_directory_errors != nil {
				return stdout_lines, ssh_directory_errors
			}

			command_escaped, command_escaped_errors := common.EscapeString(command, "'")
			if command_escaped_errors != nil {
				errors = append(errors, command_escaped_errors)
				return stdout_lines, errors
			}

			full_command := "ssh -i " + ssh_directory.GetPathAsString() + "/" + destination_user.GetFullyQualifiedUsername() + " '" + command_escaped +  "' > " + filename_stdout + " 2> " + filename_stderr + " | touch " + filename_stdout + " && touch " + filename_stderr
			fmt.Println(full_command)
			execute_unsafe_command_simple(full_command)

			if _, opening_stdout_error := os.Stat(filename_stdout); opening_stdout_error == nil {
				file_stdout, file_stdout_errors := os.Open(filename_stdout)
				if file_stdout_errors != nil {
					errors = append(errors, file_stdout_errors)
				} else {
					defer file_stdout.Close()
					stdout_scanner := bufio.NewScanner(file_stdout)
					stdout_scanner_buffer := make([]byte, maxCapacity)
					stdout_scanner.Buffer(stdout_scanner_buffer, maxCapacity)
					stdout_scanner.Split(bufio.ScanLines)
					for stdout_scanner.Scan() {
						current_text := stdout_scanner.Text()
						if current_text != "" {
							stdout_lines = append(stdout_lines, current_text)
						}
					}
				}
			}

			if _, opening_stderr_error := os.Stat(filename_stderr); opening_stderr_error == nil {
				file_stderr, file_stderr_errors := os.Open(filename_stderr)
				if file_stderr_errors != nil {
					errors = append(errors, file_stderr_errors)
				} else {
					defer file_stderr.Close()
					stderr_scanner := bufio.NewScanner(file_stderr)
					stderr_scanner_buffer := make([]byte, maxCapacity)
					stderr_scanner.Buffer(stderr_scanner_buffer, maxCapacity)
					stderr_scanner.Split(bufio.ScanLines)
					for stderr_scanner.Scan() {
						current_text := stderr_scanner.Text()
						if current_text != "" {
							errors = append(errors, fmt.Errorf("%s", current_text))
						}
					}
				}
			}

			if len(errors) > 0 {
				return stdout_lines, errors
			}

			return stdout_lines, nil
		},
		ExecuteUnsafeCommandUsingFiles: func(command string, command_data string) ([]string, []error) {
			lock.Lock()
			defer lock.Unlock()
			var errors []error
			var stdout_lines []string

			io_absolute_directory, io_absolute_directory_errors := getDirectoryIOAbsoluteDirectory()
			if io_absolute_directory_errors != nil {
				return stdout_lines, io_absolute_directory_errors
			}

			directory := io_absolute_directory.GetPathAsString()
			uuid, _ := ioutil.ReadFile("/proc/sys/kernel/random/uuid")
			time_now := time.Now().UnixNano()
			filename := directory + "/" + fmt.Sprintf("%v%s.sql", time_now, string(uuid))
			filename_stdout := directory + "/" + fmt.Sprintf("%v%s-stdout.sql", time_now, string(uuid))
			filename_stderr := directory + "/" + fmt.Sprintf("%v%s-stderr.sql", time_now, string(uuid))
			defer cleanup_files(filename, filename_stdout, filename_stderr)

			ioutil.WriteFile(filename, []byte(command_data), 0600)
			full_command := command + " < " + filename +  " > " + filename_stdout + " 2> " + filename_stderr + " | touch " + filename_stdout + " && touch " + filename_stderr
			execute_unsafe_command_simple(full_command)

			if _, opening_stdout_error := os.Stat(filename_stdout); opening_stdout_error == nil {
				file_stdout, file_stdout_errors := os.Open(filename_stdout)
				if file_stdout_errors != nil {
					errors = append(errors, file_stdout_errors)
				} else {
					defer file_stdout.Close()
					stdout_scanner := bufio.NewScanner(file_stdout)
					stdout_scanner_buffer := make([]byte, maxCapacity)
					stdout_scanner.Buffer(stdout_scanner_buffer, maxCapacity)
					stdout_scanner.Split(bufio.ScanLines)
					for stdout_scanner.Scan() {
						current_text := stdout_scanner.Text()
						if current_text != "" {
							stdout_lines = append(stdout_lines, current_text)
						}
					}
				}
			}

			if _, opening_stderr_error := os.Stat(filename_stderr); opening_stderr_error == nil {
				file_stderr, file_stderr_errors := os.Open(filename_stderr)
				if file_stderr_errors != nil {
					errors = append(errors, file_stderr_errors)
				} else {
					defer file_stderr.Close()
					stderr_scanner := bufio.NewScanner(file_stderr)
					stderr_scanner_buffer := make([]byte, maxCapacity)
					stderr_scanner.Buffer(stderr_scanner_buffer, maxCapacity)
					stderr_scanner.Split(bufio.ScanLines)
					for stderr_scanner.Scan() {
						current_text := stderr_scanner.Text()
						if current_text != "" {
							errors = append(errors, fmt.Errorf("%s", current_text))
						}
					}
				}
			}

			if len(errors) > 0 {
				return stdout_lines, errors
			}

			return stdout_lines, nil
		},
	}
	setUsername(username)

	errors := validate()

	if errors != nil {
		return nil, errors
	}

	go process_clean_up()


	return &x, nil
}

