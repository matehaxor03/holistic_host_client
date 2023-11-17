package main

import (
	host_client "github.com/matehaxor03/holistic_host_client/host_client"
	"os"
)

func main(){
	host_c, host_c_errors := host_client.NewHostClient()
	if host_c_errors != nil {
		os.Exit(1)
	}

	host_c.Validate()

	os.Exit(0)
}