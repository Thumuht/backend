package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.30

import (
	"backend/pkg/db"
	"backend/pkg/gql/graph/model"
	"backend/pkg/utils"
	"context"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// MessageID is the resolver for the messageId field.
func (r *messageResolver) MessageID(ctx context.Context, obj *db.Message) (int, error) {
	return int(obj.ID), nil
}

// CreateUser is the resolver for the createUser field.
func (r *mutationResolver) CreateUser(ctx context.Context, input model.NewUser) (*db.User, error) {
	pw, err := bcrypt.GenerateFromPassword([]byte(input.Password), 0)
	if err != nil {
		return nil, err
	}

	user := &db.User{
		LoginName: input.LoginName,
		Password:  string(pw),
	}

	_, err = r.DB.NewInsert().Model(user).Returning("*").Exec(ctx)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// FollowUser is the resolver for the followUser field.
func (r *mutationResolver) FollowUser(ctx context.Context, input int) (bool, error) {
	gctx, err := utils.GinContextFromContext(ctx)
	if err != nil {
		return false, err
	}

	_, err = r.DB.NewInsert().Model(&db.Follow{
		FollowFromId: int32(gctx.GetInt("userId")), // me
		FollowToId:   int32(input),
	}).Exec(ctx)

	if err != nil {
		return false, err
	}

	return true, nil
}

// UnfollowUser is the resolver for the unfollowUser field.
func (r *mutationResolver) UnfollowUser(ctx context.Context, input int) (bool, error) {
	meId, err := utils.GetMe(ctx)
	if err != nil {
		return false, err
	}

	_, err = r.DB.NewDelete().Model((*db.Follow)(nil)).
		Where("follow_from_id = ?", meId).
		Where("follow_to_id = ?", input).
		Exec(ctx)

	if err != nil {
		return false, err
	}

	return true, nil
}

// DeleteUser is the resolver for the deleteUser field.
func (r *mutationResolver) DeleteUser(ctx context.Context, input int) (bool, error) {
	_, err := r.DB.NewDelete().Model((*db.User)(nil)).Where("user_id = ?", input).Exec(ctx)
	if err != nil {
		return false, err
	}
	return true, nil
}

// UpdateUser is the resolver for the updateUser field.
func (r *mutationResolver) UpdateUser(ctx context.Context, input model.UpdateUser) (*db.User, error) {
	// get userId from context
	gctx, err := utils.GinContextFromContext(ctx)
	if err != nil {
		return nil, err
	}
	user := &db.User{
		ID:       int32(gctx.GetInt("userId")),
		Nickname: *input.Nickname,
		Password: *input.Password,
		Email:    *input.Email,
		About:    *input.About,
		Avatar:   *input.Avatar,
	}
	_, err = r.DB.NewUpdate().Model(user).OmitZero().Returning("*").
		Where("user_id = ?", user.ID).Exec(ctx)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// BlockUser is the resolver for the blockUser field.
func (r *mutationResolver) BlockUser(ctx context.Context, input int) (bool, error) {
	meId, err := utils.GetMe(ctx)
	if err != nil {
		return false, err
	}

	_, err = r.DB.NewInsert().Model(&db.Block{
		BlockFromId: int32(meId),
		BlockToId:   int32(input),
	}).Exec(ctx)

	if err != nil {
		return false, err
	}

	return true, nil
}

// UnblockUser is the resolver for the unblockUser field.
func (r *mutationResolver) UnblockUser(ctx context.Context, input int) (bool, error) {
	meId, err := utils.GetMe(ctx)
	if err != nil {
		return false, err
	}

	_, err = r.DB.NewDelete().Model((*db.Block)(nil)).
		Where("block_from_id = ?", meId).
		Where("block_to_id = ?", input).
		Exec(ctx)

	if err != nil {
		return false, err
	}

	return true, nil
}

// SendMessage is the resolver for the sendMessage field.
func (r *mutationResolver) SendMessage(ctx context.Context, input model.MessageInput) (bool, error) {
	meId, err := utils.GetMe(ctx)
	if err != nil {
		return false, err
	}

	_, err = r.DB.NewInsert().Model(&db.Message{
		UserFrom: int32(meId),
		UserTo:   int32(input.ToID),
		Content:  input.Content,
	}).Exec(ctx)

	if err != nil {
		return false, err
	}

	return true, nil
}

// Login is the resolver for the login field.
func (r *mutationResolver) Login(ctx context.Context, input model.LoginSession) (string, error) {
	var user db.User

	err := r.DB.NewSelect().Model(&user).Where("login_name = ?", input.LoginName).Scan(ctx)
	if err != nil {
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return "", fmt.Errorf("login %s %s, input pw %s %w", input.LoginName, user.Password, input.Password, err)
	}

	token := utils.GenToken()
	r.Cache.Sessions.Set(token, int(user.ID))

	return token, nil
}

// Logout is the resolver for the logout field.
func (r *mutationResolver) Logout(ctx context.Context) (bool, error) {
	gctx, err := utils.GinContextFromContext(ctx)
	if err != nil {
		return false, fmt.Errorf("cannot get gin context, access denied: %w", err)
	}

	token := gctx.GetHeader("Token")
	return true, r.Cache.Sessions.Remove(token)
}

// Users is the resolver for the users field.
func (r *queryResolver) Users(ctx context.Context, input model.GetUserInput) ([]db.User, error) {
	var users []db.User
	err := r.DB.NewSelect().Model(&users).Relation("Post").Relation("Comment").
		Order(input.OrderBy.String() + " " + input.Order.String()).Limit(input.Limit).
		Offset(input.Offset).
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	return users, nil
}

// GetUserByID is the resolver for the getUserById field.
func (r *queryResolver) GetUserByID(ctx context.Context, input int) (*db.User, error) {
	var user db.User
	err := r.DB.NewSelect().Model(&user).Where("user_id = ?", input).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Me is the resolver for the me field.
func (r *queryResolver) Me(ctx context.Context) (*db.User, error) {
	meId, err := utils.GetMe(ctx)
	if err != nil {
		return nil, err
	}

	var user db.User
	err = r.DB.NewSelect().Model(&user).Where("user_id = ?", meId).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Follow is the resolver for the follow field.
func (r *userResolver) Follow(ctx context.Context, obj *db.User) ([]*db.User, error) {
	// find all user that obj follows
	var users []*db.User
	err := r.DB.NewSelect().Model(&users).Relation("Post").Relation("Comment").
		Where("user_id IN (SELECT follow_id FROM follow WHERE user_id = ?)", obj.ID).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return users, nil
}

// Follower is the resolver for the follower field.
func (r *userResolver) Follower(ctx context.Context, obj *db.User) ([]*db.User, error) {
	// follow all user that follows obj
	var users []*db.User
	err := r.DB.NewSelect().Model(&users).Relation("Post").Relation("Comment").
		Where("user_id IN (SELECT user_id FROM follow WHERE follow_id = ?)", obj.ID).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return users, nil
}

// Block is the resolver for the block field.
func (r *userResolver) Block(ctx context.Context, obj *db.User) ([]*db.User, error) {
	// find all user that obj blocks
	var users []*db.User
	err := r.DB.NewSelect().Model(&users).Relation("Post").Relation("Comment").
		Where("user_id IN (SELECT block_to_id FROM block WHERE block_from = ?)", obj.ID).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return users, nil
}

// Message returns MessageResolver implementation.
func (r *Resolver) Message() MessageResolver { return &messageResolver{r} }

// User returns UserResolver implementation.
func (r *Resolver) User() UserResolver { return &userResolver{r} }

type messageResolver struct{ *Resolver }
type userResolver struct{ *Resolver }
