package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.26

import (
	"backend/pkg/db"
	"backend/pkg/gql/graph/model"
	"context"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/spf13/viper"
)

// ParentType is the resolver for the parentType field.
func (r *attachmentResolver) ParentType(ctx context.Context, obj *db.Attachment) (model.AttachmentParent, error) {
	return model.AttachmentParent(obj.ParentType), nil
}

// FileUpload is the resolver for the fileUpload field.
func (r *mutationResolver) FileUpload(ctx context.Context, input *model.PostUpload) (bool, error) {
	content, err := io.ReadAll(input.Upload.File)
	if err != nil {
		return false, err
	}

	attach := &db.Attachment{
		FileName:   input.Upload.Filename,
		ParentID:   int32(input.ParentID),
		ParentType: input.ParentType.String(),
	}
	_, err = r.DB.NewInsert().Model(attach).Exec(ctx)
	if err != nil {
		return false, err
	}
	fmt.Printf("[GQL] write file %s\n", path.Join(viper.GetString("fs_route"), input.Upload.Filename))
	os.WriteFile(path.Join(viper.GetString("fs_route"), input.Upload.Filename), content, 0666)

	return true, nil
}

// Attachment returns AttachmentResolver implementation.
func (r *Resolver) Attachment() AttachmentResolver { return &attachmentResolver{r} }

type attachmentResolver struct{ *Resolver }
