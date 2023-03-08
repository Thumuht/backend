/*
eXecutable programs for thumuht

# main

provides http backend server for thumuht

Usage:

	main

Configuration:

The program need a config.yml file in the working directory to run properly.
Settings below are required.

fs_route

	directory where saves user-upload files

addr

	socket addr where serves the service
*/
package main

import (
	"backend/pkg/config"
	"backend/pkg/forum"
)

const (
	forumAdd = "127.0.0.1:8899"
)

func main() {
	err := config.ConfigProject()
	if err != nil {
		panic("no config")
	}
	forumapp := forum.NewForum()
	forumapp.RunForum(forumAdd)
}
