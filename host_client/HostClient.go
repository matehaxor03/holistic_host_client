package host_client

import (
)

type HostClient struct {
	Validate func() []error
	CreateHostUser func(username string) (*HostUser, []error)
	DeleteHostUser func(username string) ([]error)
	Ramdisk func(disk_name string,block_size uint64) (*Ramdisk, []error)
}

func NewHostClient() (*HostClient, []error) {
	
	createHostUser := func(username string) (*HostUser, []error) {
		host_user, host_user_errors := newHostUser(username)
		
		if host_user_errors != nil {
			return nil, host_user_errors
		}

		create_errors := host_user.Create()
		
		if create_errors != nil {
			return nil, create_errors
		}

		return host_user, nil
	}

	deleteHostUser := func(username string) ([]error) {
		host_user, host_user_errors := newHostUser(username)
		
		if host_user_errors != nil {
			return host_user_errors
		}

		delete_errors := host_user.Delete()
		
		if delete_errors != nil {
			return delete_errors
		}

		return nil
	}

	ramdisk := func(disk_name string, block_size uint64) (*Ramdisk, []error) {
		ramdisk, ramdisk_errors := newRamdisk(disk_name, block_size)
		if ramdisk_errors != nil {
			return nil, ramdisk_errors
		}

		return ramdisk, nil
	}

	validate := func() []error {
		return nil
	}

	x := HostClient{
		Validate: func() []error {
			return validate()
		},
		CreateHostUser: func(username string) (*HostUser, []error) {
			return createHostUser(username)
		},
		DeleteHostUser: func(username string) ([]error) {
			return deleteHostUser(username)
		},
		Ramdisk: func(disk_name string, block_size uint64) (*Ramdisk, []error) {
			return ramdisk(disk_name, block_size)
		},
	}
	//setHostClient(&x)

	errors := validate()

	if errors != nil {
		return nil, errors
	}

	return &x, nil
}

