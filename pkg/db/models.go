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

	Post         []*Post         `json:"post" bun:"rel:has-many,join:user_id=post_userid"`
	Comment      []*Comment      `json:"comment" bun:"rel:has-many,join:user_id=comment_userid"`
	Follow       []*Follow       `json:"follow" bun:"rel:has-many,join:user_id=follow_from"`
	Follower     []*Follow       `json:"follower" bun:"rel:has-many,join:user_id=follow_to"`
	BookmarkList []*BookmarkList `json:"userBookmarkList" bun:"rel:has-many,join:user_id=bookmark_list_userid"`
}

type Post struct {
	bun.BaseModel `bun:"table:post"`

	ID        int32     `json:"id" bun:"post_id,pk,autoincrement"`
	Title     string    `json:"title" bun:"title"`
	Content   string    `json:"content" bun:"content"`
	Position  string    `json:"position" bun:"position"`
	View      int32     `json:"view" bun:"view"`
	Like      int32     `json:"like" bun:"like"`
	CommentsNum int32   `json:"commentsNum" bun:"comments_num"`
	UserID    int32     `json:"userId" bun:"post_userid"`
	CreatedAt time.Time `json:"createdAt" bun:",nullzero,notnull,default:current_timestamp"`
	UpdatedAt time.Time `json:"updatedAt" bun:",nullzero,notnull,default:current_timestamp"`

	User       *User         `json:"user" bun:"rel:belongs-to,join:post_userid=user_id,on_delete:cascade"`
	Comment    []*Comment    `json:"comment" bun:"rel:has-many,join:post_id=comment_postid"`
	Attachment []*Attachment `json:"attachment" bun:"rel:has-many,join:post_id=attachment_parentid,join:type=parent_type,polymorphic"`
	Tag []*Tag `json:"tag" bun:"m2m:post_tag,join:Tag=Post"`
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

// TODO(wj): follow, Message
type Follow struct {
	bun.BaseModel `bun:"table:follow"`

	FollowFromId int32 `bun:"follow_from,pk"`
	FollowToId   int32 `bun:"follow_to,pk"`

	FollowFrom *User `json:"followFrom" bun:"rel:belongs-to,join:follow_from=user_id,on_delete:cascade"`
	FollowTo   *User `json:"followTo" bun:"rel:belongs-to,join:follow_to=user_id,on_delete:cascade"`

	CreatedAt int32 `bun:",nullzero,notnull,default:current_timestamp"`
}

type Message struct {
	bun.BaseModel `bun:"table:message"`

	ID       int32  `json:"messageId" bun:"message_id,pk"`
	UserFrom int32  `json:"userFrom" bun:"user_from"`
	UserTo   int32  `json:"userTo" bun:"user_to"`
	Content  string `json:"content" bun:"content"`
	CreateAt int32  `json:"createAt" bun:"create_at"`
}

type Bookmark struct {
	bun.BaseModel `bun:"table:bookmark"`

	PostID         int32         `bun:"bookmark_postid,pk"`
	Post           *Post         `json:"post" bun:"rel:belongs-to,join:bookmark_postid=post_id"`
	BookmarkListID int32         `bun:"bookmark_bookmarklistid,pk"`
	BookmarkList   *BookmarkList `json:"bookmarkList" bun:"rel:belongs-to,join:bookmark_bookmarklistid=bookmark_list_id"`
}

type BookmarkList struct {
	bun.BaseModel `bun:"table:bookmark_list"`

	ID        int32     `json:"id" bun:"bookmark_list_id,pk"`
	List      string    `json:"list" bun:"bookmark_list"`
	UserID    int32     `bun:"bookmark_list_userid"`
	CreatedAt time.Time `json:"createdAt" bun:",nullzero,notnull,default:current_timestamp"`
	UpdatedAt time.Time `json:"updatedAt" bun:",nullzero,notnull,default:current_timestamp"`

	User *User `json:"user" bun:"rel:belongs-to,join:bookmark_list_userid=user_id,on_delete:cascade"`
	// many-to-many through Bookmark
	Post []*Post `json:"post" bun:"m2m:bookmark,join:Post=BookmarkList"`
}

type Block struct {
	bun.BaseModel `bun:"table:block"`

	BlockFromId int32 `bun:"block_from,pk"`
	BlockToId   int32 `bun:"block_to,pk"`

	BlockFrom *User `json:"blockFrom" bun:"rel:belongs-to,join:block_from=user_id,on_delete:cascade"`
	BlockTo   *User `json:"blockTo" bun:"rel:belongs-to,join:block_to=user_id,on_delete:cascade"`

	CreatedAt int32 `bun:",nullzero,notnull,default:current_timestamp"`
}

type Tag struct {
	bun.BaseModel `bun:"table:tag"`

	TagId   int32  `json:"tagId" bun:"tag_id,pk"`
	TagName string `json:"tagName" bun:"tag_name"`

	Post []*Post `json:"post" bun:"m2m:post_tag,join:Tag=Post"`
}

type PostTag struct {
	bun.BaseModel `bun:"table:post_tag"`

	PostID int32 `bun:"post_tag_postid,pk"`
	TagID  int32 `bun:"post_tag_tagid,pk"`

	Post *Post `json:"post" bun:"rel:belongs-to,join:post_tag_postid=post_id,on_delete:cascade"`
	Tag  *Tag  `json:"tag" bun:"rel:belongs-to,join:post_tag_tagid=tag_id,on_delete:cascade"`
}
