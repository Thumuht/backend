package main

import "backend/pkg/forum"

const (
	forumAdd = "127.0.0.1:8899"
)

func main() {
	forumapp := forum.NewForum()
	forumapp.RunForum(forumAdd)
}
