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
