package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.31

import (
	"backend/pkg/db"
	"backend/pkg/gql/graph/model"
	"context"
	"time"
)

// CreateComment is the resolver for the createComment field.
func (r *mutationResolver) CreateComment(ctx context.Context, input model.NewComment) (*db.Comment, error) {
	comment := &db.Comment{
		UserID:  int32(input.UserID),
		PostID:  int32(input.PostID),
		Content: *input.Content,
	}
	_, err := r.DB.NewInsert().Model(comment).Returning("*").Exec(ctx)
	if err != nil {
		return nil, err
	}
	return comment, nil
}

// DeleteComment is the resolver for the deleteComment field.
func (r *mutationResolver) DeleteComment(ctx context.Context, input int) (bool, error) {
	_, err := r.DB.NewDelete().Model((*db.Comment)(nil)).Where("comment_id = ?", input).Exec(ctx)
	if err != nil {
		return false, err
	}
	return true, nil
}

// UpdateComment is the resolver for the updateComment field.
func (r *mutationResolver) UpdateComment(ctx context.Context, input model.UpdateComment) (*db.Comment, error) {
	comment := &db.Comment{
		ID:        int32(input.CommentID),
		UpdatedAt: time.Now(),
	}

	if input.Content != nil {
		comment.Content = *input.Content
	}
	_, err := r.DB.NewUpdate().Model(comment).OmitZero().WherePK().Returning("*").Exec(ctx)
	if err != nil {
		return nil, err
	}
	return comment, nil
}

// LikeComment is the resolver for the likeComment field.
func (r *mutationResolver) LikeComment(ctx context.Context, input int) (bool, error) {
	// first see cache
	// if not exists, fetch from db and write to cache
	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	if like, ok := r.Cache.CommentLike.Get(input); ok {
		r.Cache.CommentLike.Set(input, *like+1)
	} else {
		var comment db.Comment
		err := r.DB.NewSelect().Model(&comment).Where("comment_id = ?", input).Scan(ctx)
		if err != nil {
			return false, err
		}
		r.Cache.CommentLike.Set(input, int(comment.Like)+1)
	}

	return true, nil
}

// DislikeComment is the resolver for the dislikeComment field.
func (r *mutationResolver) DislikeComment(ctx context.Context, input int) (bool, error) {
	// first see cache
	// if not exists, fetch from db and write to cache
	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	if like, ok := r.Cache.CommentLike.Get(input); ok {
		r.Cache.CommentLike.Set(input, *like-1)
	} else {
		var comment db.Comment
		err := r.DB.NewSelect().Model(&comment).Where("comment_id = ?", input).Scan(ctx)
		if err != nil {
			return false, err
		}
		r.Cache.CommentLike.Set(input, int(comment.Like)-1)
	}

	return true, nil
}

// Comment is the resolver for the comment field.
func (r *queryResolver) Comment(ctx context.Context, input model.GetCommentInput) ([]db.Comment, error) {
	var comments []db.Comment
	err := r.DB.NewSelect().Model(&comments).Relation("Post").Relation("User").Relation("Attachment").
		Order("comment." + input.OrderBy.String() + " " + input.Order.String()).Limit(input.Limit).
		Offset(input.Offset).Scan(ctx)
	if err != nil {
		return nil, err
	}

	for i := range comments {
		if like, ok := r.Cache.PostLike.Get(int(comments[i].ID)); ok {
			comments[i].Like = int32(*like)
		}
	}

	return comments, nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
