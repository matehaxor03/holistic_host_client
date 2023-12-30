package host_client

import (
	validate "github.com/matehaxor03/holistic_validator/validate"
	common "github.com/matehaxor03/holistic_common/common"
	"fmt"
	"strconv"
)

type Ramdisk struct {
	Validate func() []error
	Create func() []error
}

func newRamdisk(disk_name string, block_size uint64) (*Ramdisk, []error) {
	bashCommand := common.NewBashCommand()
	verify := validate.NewValidator()
	var this_disk_name string
	var this_block_size uint64

	setDiskName := func(disk_name string) {
		this_disk_name = disk_name
	}

	getDiskName := func() string {
		return this_disk_name
	}

	setBlockSize := func(block_size uint64) {
		this_block_size = block_size
	}

	getBlockSize := func() uint64 {
		return this_block_size
	}

	validate := func() []error {
		var errors []error
		temp_disk_name := getDiskName()
		temp_block_size := getBlockSize()

		disk_name_errors := verify.ValidateDirectoryName(temp_disk_name)
		if disk_name_errors != nil {
			errors = append(errors, disk_name_errors...)
		}

		if temp_block_size == 0 {
			errors = append(errors, fmt.Errorf("ram disk cannot have 0 blocks"))
		}

		if len(errors) > 0 {
			return errors
		}
		return nil
	}

	create := func() []error {
		shell_command := "diskutil erasevolume HFS+ \"" + getDiskName() + "\" `hdiutil attach -nomount ram://" + strconv.FormatUint(getBlockSize(), 10)
		_, std_error := bashCommand.ExecuteUnsafeCommandUsingFilesWithoutInputFile(shell_command)
		if std_error != nil {
			return std_error
		}
		return nil
	}

	x := Ramdisk{
		Validate: func() []error {
			return validate()
		},
		Create: func() []error {
			return create()
		},
	}
	setDiskName(disk_name)
	setBlockSize(block_size)

	errors := validate()

	if errors != nil {
		return nil, errors
	}

	return &x, nil
}

