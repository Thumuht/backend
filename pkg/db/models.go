package db

import "github.com/uptrace/bun"

type User struct {
	bun.BaseModel `bun:"table:user"`

	ID        int32  `json:"id" bun:"user_id,pk,autoincrement"`
	LoginName string `json:"loginName" bun:"login_name,unique,notnull"`
	Nickname  string `json:"nickname" bun:"nickname"`
	Password  string `json:"password" bun:"pwhash"`
	Email     string `json:"email" bun:"email"`
	About     string `json:"about" bun:"about"`

	Post    []*Post    `json:"post" bun:"rel:has-many,join:user_id=post_userid"`
	Comment []*Comment `json:"comment" bun:"rel:has-many,join:user_id=comment_userid"`
}

type Post struct {
	bun.BaseModel `bun:"table:post"`

	ID      int32  `json:"id" bun:"post_id,pk,autoincrement"`
	Title   string `json:"title" bun:"title"`
	Content string `json:"content" bun:"content"`
	UserID  int32  `json:"userId" bun:"post_userid"`

	User    *User      `json:"user" bun:"rel:belongs-to,join:post_userid=user_id"`
	Comment []*Comment `json:"comment" bun:"rel:has-many,join:post_id=comment_postid"`
}

type Comment struct {
	bun.BaseModel `bun:"table:comment"`

	ID      int32  `json:"id" bun:"comment_id,pk,autoincrement"`
	Content string `json:"content" bun:"content"`
	UserID  int32  `json:"userId" bun:"comment_userid"`
	PostID  int32  `json:"postId" bun:"comment_postid"`

	User *User `json:"user" bun:"rel:belongs-to,join:comment_userid=user_id"`
	Post *Post `json:"post" bun:"rel:belongs-to,join:comment_postid=post_id"`
}
