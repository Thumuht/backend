package graph

import (
	"backend/pkg/db"
	"context"
)

func (r *mutationResolver) sendMsgTo(ctx context.Context, msg string, from int, to int) error {
	message := &db.Message{
		UserFrom: int32(from),
		UserTo:   int32(to),
		Content:  msg,
		IsNew:    true,
	}

	if ch, ok := r.Cache.Notifier.Get(int(message.UserTo)); ok {
		go func() {
			*ch <- message
		}()
		message.IsNew = false
	}

	_, err := r.DB.NewInsert().Model(message).Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}
