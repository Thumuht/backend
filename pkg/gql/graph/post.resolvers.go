package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.31

import (
	"backend/pkg/db"
	"backend/pkg/gql/graph/model"
	"backend/pkg/utils"
	"context"
	"fmt"
	"sort"
	"time"
)

// CreatePost is the resolver for the createPost field.
func (r *mutationResolver) CreatePost(ctx context.Context, input model.NewPost) (*db.Post, error) {
	post := &db.Post{
		UserID: int32(input.UserID),
	}
	if input.Content != nil {
		post.Content = *input.Content
	}

	if input.Title != nil {
		post.Title = *input.Title
	}

	if input.Tag != nil {
		post.Tag = *input.Tag
	}

	if input.Position != nil {
		post.Position = *input.Position
	}

	_, err := r.DB.NewInsert().Model(post).Returning("*").Exec(ctx)
	if err != nil {
		return nil, err
	}

	// inform all followers of the new post
	var followers []db.Follow
	err = r.DB.NewSelect().Model(&followers).Where("follow_to_id = ?", input.UserID).Scan(ctx)
	if err != nil {
		return nil, err
	}

	msg := fmt.Sprintf("你关注的用户发布了新的帖子: %s", *input.Title)
	for _, follower := range followers {
		if r.sendMsgTo(ctx, msg, utils.SysAccountID, int(follower.FollowFromId)) != nil {
			return nil, err
		}
	}

	return post, nil
}

// UpdatePost is the resolver for the updatePost field.
func (r *mutationResolver) UpdatePost(ctx context.Context, input model.UpdatePost) (*db.Post, error) {
	post := &db.Post{
		ID:        int32(input.PostID),
		UpdatedAt: time.Now(),
	}

	if input.Content != nil {
		post.Content = *input.Content
	}

	if input.Title != nil {
		post.Title = *input.Title
	}

	_, err := r.DB.NewUpdate().Model(post).OmitZero().WherePK().Returning("*").Exec(ctx)
	if err != nil {
		return nil, err
	}
	return post, nil
}

// DeletePost is the resolver for the deletePost field.
func (r *mutationResolver) DeletePost(ctx context.Context, postID int) (bool, error) {
	_, err := r.DB.NewDelete().Model((*db.Post)(nil)).Where("post_id = ?", postID).Exec(ctx)
	if err != nil {
		return false, err
	}
	return true, nil
}

// MarkPost is the resolver for the markPost field.
func (r *mutationResolver) MarkPost(ctx context.Context, input int) (bool, error) {
	userTok, _ := utils.GetMe(ctx)
	userId, _ := r.Cache.Sessions.Get(userTok)
	markPost := &db.Bookmark{
		BookmarkPostID: int32(input),
		BookmarkUserID: int32(*userId),
	}
	_, err := r.DB.NewInsert().Model(markPost).Returning("*").Exec(ctx)
	if err != nil {
		return false, err
	}
	return true, nil
}

// UnmarkPost is the resolver for the unmarkPost field.
func (r *mutationResolver) UnmarkPost(ctx context.Context, input int) (bool, error) {
	userTok, _ := utils.GetMe(ctx)

	userId, _ := r.Cache.Sessions.Get(userTok)

	_, err := r.DB.NewDelete().Model((*db.Bookmark)(nil)).Where("bookmark_post_id = ? AND bookmark_user_id = ?", input, *userId).Exec(ctx)
	if err != nil {
		return false, err
	}
	return true, nil
}

// LikePost is the resolver for the likePost field.
func (r *mutationResolver) LikePost(ctx context.Context, input int) (bool, error) {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()
	var post db.Post
	r.DB.NewSelect().Model(&post).Where("post_id = ?", input).Scan(ctx)
	r.DB.NewUpdate().Model(&post).Set("like = ?", post.Like+1).WherePK().Exec(ctx)
	r.sendMsgTo(ctx, "one user liked your post!", utils.SysAccountID, int(post.UserID))
	return true, nil
}

// DislikePost is the resolver for the dislikePost field.
func (r *mutationResolver) DislikePost(ctx context.Context, input int) (bool, error) {
	// first see cache
	// if not exists, fetch from db and write to cache
	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	var post db.Post
	r.DB.NewSelect().Model(&post).Where("post_id = ?", input).Scan(ctx)
	r.DB.NewUpdate().Model(&post).Set("like = ?", post.Like-1).WherePK().Exec(ctx)
	return true, nil
}

// Posts is the resolver for the posts field.
func (r *queryResolver) Posts(ctx context.Context, input model.GetPostInput) ([]*db.Post, error) {
	var posts []*db.Post
	// get me if login
	meTok, err := utils.GetMe(ctx)
	if err != nil {
		return nil, err
	}

	meId, _ := r.Cache.Sessions.Get(meTok)

	var pquery = r.DB.NewSelect().Model(&posts).Relation("User").Relation("Comment").Relation("Attachment").
		Order("post." + input.OrderBy.String() + " " + input.Order.String()).Offset(input.Offset).Limit(input.Limit)

	if meId != nil {
		pquery = pquery.Where("post_userid NOT IN (SELECT block_to_id FROM block WHERE block_from_id = ?)", *meId)
	}

	if *input.Followed {
		pquery = pquery.Where("post_userid IN (SELECT follow_to_id FROM follow WHERE follow_from_id = ?)", *meId)
	}

	if input.Tags != nil {
		pquery = pquery.Where("tag = ?", input.Tags)
	}

	err = pquery.Scan(ctx)
	if err != nil {
		return nil, err
	}

	// range is a copy..

	// sort by like in cache, or view in cache
	if input.OrderBy == model.PostOrderByLike {
		// sort by like
		sort.Slice(posts, func(i, j int) bool {
			return posts[i].Like > posts[j].Like
		})
	} else if input.OrderBy == model.PostOrderByView {
		// sort by view
		sort.Slice(posts, func(i, j int) bool {
			return posts[i].View > posts[j].View
		})
	}

	return posts, nil
}

// PostDetail is the resolver for the postDetail field.
func (r *queryResolver) PostDetail(ctx context.Context, input int) (*db.Post, error) {
	var post db.Post

	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	err := r.DB.NewSelect().Model(&post).Relation("Comment").Relation("Attachment").Relation("User").Where("post_id = ?", input).Scan(ctx)
	if err != nil {
		return nil, err
	}

	// fetch from db, +1, update
	r.DB.NewUpdate().Model(&post).Set("view = ?", post.View+1).WherePK().Exec(ctx)

	post.CommentsNum = int32(len(post.Comment))
	return &post, nil
}
