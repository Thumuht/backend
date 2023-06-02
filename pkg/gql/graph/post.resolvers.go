package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.31

import (
	"backend/pkg/db"
	"backend/pkg/gql/graph/model"
	"backend/pkg/utils"
	"context"
	"sort"
	"time"
)

// CreatePost is the resolver for the createPost field.
func (r *mutationResolver) CreatePost(ctx context.Context, input model.NewPost) (*db.Post, error) {
	post := &db.Post{
		UserID:  int32(input.UserID),
		Content: *input.Content,
		Title:   *input.Title,
	}
	_, err := r.DB.NewInsert().Model(post).Returning("*").Exec(ctx)
	if err != nil {
		return nil, err
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

// NewBookmarkList is the resolver for the newBookmarkList field.
func (r *mutationResolver) NewBookmarkList(ctx context.Context, input string) (*db.BookmarkList, error) {
	gctx, err := utils.GinContextFromContext(ctx)
	if err != nil {
		return nil, err
	}
	cTime := time.Now()
	bookmarkList := &db.BookmarkList{
		List:      input,
		UserID:    int32(gctx.GetInt("userId")),
		CreatedAt: cTime,
		UpdatedAt: cTime,
	}
	_, err = r.DB.NewInsert().Model(bookmarkList).Returning("*").Exec(ctx)
	if err != nil {
		return nil, err
	}
	return bookmarkList, nil
}

// MarkPost is the resolver for the markPost field.
func (r *mutationResolver) MarkPost(ctx context.Context, input model.NewMarkPost) (bool, error) {
	markPost := &db.Bookmark{
		PostID:         int32(input.PostID),
		BookmarkListID: int32(input.BookmarkListID),
	}
	_, err := r.DB.NewInsert().Model(markPost).Returning("*").Exec(ctx)
	if err != nil {
		return false, err
	}
	return true, nil
}

// UnmarkPost is the resolver for the unmarkPost field.
func (r *mutationResolver) UnmarkPost(ctx context.Context, input model.NewMarkPost) (bool, error) {
	_, err := r.DB.NewDelete().Model((*db.Bookmark)(nil)).Where("post_id = ? AND bookmark_list_id = ?", input.PostID, input.BookmarkListID).Exec(ctx)
	if err != nil {
		return false, err
	}
	return true, nil
}

// LikePost is the resolver for the likePost field.
func (r *mutationResolver) LikePost(ctx context.Context, input int) (bool, error) {
	// first see cache
	// if not exists, fetch from db and write to cache
	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	if like, ok := r.Cache.PostLike.Get(input); ok {
		r.Cache.PostLike.Set(input, *like+1)
	} else {
		var post db.Post
		err := r.DB.NewSelect().Model(&post).Where("post_id = ?", input).Scan(ctx)
		if err != nil {
			return false, err
		}
		r.Cache.PostLike.Set(input, int(post.Like)+1)
	}

	return true, nil
}

// DislikePost is the resolver for the dislikePost field.
func (r *mutationResolver) DislikePost(ctx context.Context, input int) (bool, error) {
	// first see cache
	// if not exists, fetch from db and write to cache
	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	if like, ok := r.Cache.PostLike.Get(input); ok {
		r.Cache.PostLike.Set(input, *like-1)
	} else {
		var post db.Post
		err := r.DB.NewSelect().Model(&post).Where("post_id = ?", input).Scan(ctx)
		if err != nil {
			return false, err
		}
		r.Cache.PostLike.Set(input, int(post.Like)-1)
	}

	return true, nil
}

// Posts is the resolver for the posts field.
func (r *queryResolver) Posts(ctx context.Context, input model.GetPostInput) ([]*db.Post, error) {
	var posts []*db.Post
	// get me if login
	meId, err := utils.GetMe(ctx)
	if err != nil {
		return nil, err
	}

	var pquery = r.DB.NewSelect().Model(&posts).Relation("User").Relation("Comment").Relation("Attachment").
		Order("post." + input.OrderBy.String()+" "+input.Order.String()).Offset(input.Offset).Limit(input.Limit).
		Where("post_userid NOT IN (SELECT block_to_id FROM block WHERE block_from_id = ?)", meId)

	if *input.Followed {
		pquery = pquery.Where("post_userid IN (SELECT follow_to_id FROM follow WHERE follow_from_id = ?)", meId)
	}

	if !*input.All {
		// tag in input.Tags
		// join post, post_tag, tag
		pquery = pquery.Join("JOIN post_tag ON post_id = post_tag_postid").
			Join("JOIN tag ON post_tag_tagid = tag_id").
			Where("tag_name IN ?", input.Tags)
	}

	err = pquery.Scan(ctx)
	if err != nil {
		return nil, err
	}

	// range is a copy..
	for i := range posts {
		if like, ok := r.Cache.PostLike.Get(int(posts[i].ID)); ok {
			posts[i].Like = int32(*like)
		}
		if view, ok := r.Cache.PostView.Get(int(posts[i].ID)); ok {
			posts[i].View = int32(*view)
		}
	}

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

	// first see cache
	// if not exists, fetch from db and write to cache
	if view, ok := r.Cache.PostView.Get(input); ok {
		post.View = int32(*view) + 1
	} else {
		post.View += 1
	}

	if like, ok := r.Cache.PostLike.Get(input); ok {
		post.Like = int32(*like)
	}

	r.Cache.PostView.Set(input, int(post.View))
	post.CommentsNum = int32(len(post.Comment))
	return &post, nil
}
