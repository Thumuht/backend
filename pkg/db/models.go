package db

import (
	"time"

	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bun:"table:user"`

	ID        int32  `json:"id" bun:"user_id,pk,autoincrement"`
	LoginName string `json:"loginName" bun:"login_name,unique,notnull"`
	Nickname  string `json:"nickname" bun:"nickname"`
	Password  string `json:"password" bun:"pwhash"`
	Email     string `json:"email" bun:"email"`
	About     string `json:"about" bun:"about"`
	Avatar    string `json:"avatar" bun:"avatar"`

	Post         []*Post    `json:"post" bun:"rel:has-many,join:user_id=post_userid"`
	Comment      []*Comment `json:"comment" bun:"rel:has-many,join:user_id=comment_userid"`
	Follow       []*Follow  `json:"follow" bun:"rel:has-many,join:user_id=follow_from_id"`
	Follower     []*Follow  `json:"follower" bun:"rel:has-many,join:user_id=follow_to_id"`
	Block        []*Block   `json:"block"`
	BookmarkList []*Post    `json:"userBookmarkList" bun:"m2m:bookmark,join:User=Post"`
}

type Post struct {
	bun.BaseModel `bun:"table:post"`

	ID          int32     `json:"id" bun:"post_id,pk,autoincrement"`
	Title       string    `json:"title" bun:"title"`
	Content     string    `json:"content" bun:"content"`
	Position    string    `json:"position" bun:"position"`
	View        int32     `json:"view" bun:"view"`
	Like        int32     `json:"like" bun:"like"`
	CommentsNum int32     `json:"commentsNum" bun:"comments_num"`
	UserID      int32     `json:"userId" bun:"post_userid"`
	Tag         string    `json:"tag" bun:"tag"`
	CreatedAt   time.Time `json:"createdAt" bun:",nullzero,notnull,default:current_timestamp"`
	UpdatedAt   time.Time `json:"updatedAt" bun:",nullzero,notnull,default:current_timestamp"`

	User       *User         `json:"user" bun:"rel:belongs-to,join:post_userid=user_id,on_delete:cascade"`
	Comment    []*Comment    `json:"comment" bun:"rel:has-many,join:post_id=comment_postid"`
	Attachment []*Attachment `json:"attachment" bun:"rel:has-many,join:post_id=attachment_parentid,join:type=parent_type,polymorphic"`
}

type Comment struct {
	bun.BaseModel `bun:"table:comment"`

	ID        int32     `json:"id" bun:"comment_id,pk,autoincrement"`
	Content   string    `json:"content" bun:"content"`
	Like      int32     `json:"like" bun:"like"`
	UserID    int32     `json:"userId" bun:"comment_userid"`
	PostID    int32     `json:"postId" bun:"comment_postid"`
	CreatedAt time.Time `json:"createdAt" bun:",nullzero,notnull,default:current_timestamp"`
	UpdatedAt time.Time `json:"updatedAt" bun:",nullzero,notnull,default:current_timestamp"`

	User       *User         `json:"user" bun:"rel:belongs-to,join:comment_userid=user_id,on_delete:cascade"`
	Post       *Post         `json:"post" bun:"rel:belongs-to,join:comment_postid=post_id,on_delete:cascade"`
	Attachment []*Attachment `json:"attachment" bun:"rel:has-many,join:comment_id=attachment_parentid,join:type=parent_type,polymorphic"`
}

type Attachment struct {
	bun.BaseModel `bun:"table:attachment"`

	ID         int32     `json:"id" bun:"attachment_id,pk,autoincrement"`
	ParentID   int32     `json:"parentId" bun:"attachment_parentid"`
	ParentType string    `json:"parentType" bun:"parent_type"`
	FileName   string    `json:"fileName" bun:"file_name"`
	CreatedAt  time.Time `json:"createdAt" bun:",nullzero,notnull,default:current_timestamp"`
}

type Follow struct {
	bun.BaseModel `bun:"table:follow"`

	FollowFromId int32 `bun:",pk"`
	FollowToId   int32 `bun:",pk"`

	CreatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
}

type Message struct {
	bun.BaseModel `bun:"table:message"`

	ID        int32     `json:"messageId" bun:"message_id,pk,autoincrement"`
	UserFrom  int32     `json:"userFrom" bun:"user_from"`
	UserTo    int32     `json:"userTo" bun:"user_to"`
	Content   string    `json:"content" bun:"content"`
	IsNew     bool      `json:"isNew" bun:"is_new"`
	CreatedAt time.Time `json:"createdAt" bun:"created_at,nullzero,notnull,default:current_timestamp"`
}

type Bookmark struct {
	bun.BaseModel `bun:"table:bookmark"`

	BookmarkPostID int32 `bun:",pk"`
	BookmarkUserID int32 `bun:",pk"`
	Post           *Post `bun:"rel:belongs-to,join:bookmark_post_id=post_id,on_delete:cascade"`
	User           *User `bun:"rel:belongs-to,join:bookmark_user_id=user_id,on_delete:cascade"`
}

type Block struct {
	bun.BaseModel `bun:"table:block"`

	BlockFromId int32 `bun:"block_from_id,pk"`
	BlockToId   int32 `bun:"block_to_id,pk"`
}
