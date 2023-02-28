package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.24

import (
	"backend/pkg/db"
	"backend/pkg/gql/graph/model"
	"backend/pkg/utils"
	"context"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

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
	r.Sessions[user.LoginName] = token

	return token, nil
}