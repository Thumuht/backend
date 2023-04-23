/*
Package db defines model of thumuht, and sets up the database.

TODO(wj, mid): add more test cases
*/
package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
	"github.com/uptrace/bun/extra/bundebug"
)

// init SQLiteDB for debug/test purpose only
// 'cause it stores its data on memory, and is volatile
// therefore, it's convinent for testing, but unacceptable in real use.
func InitSQLiteDB() (*bun.DB, error) {
	sqldb, err := sql.Open(sqliteshim.ShimName, "file::memory:?cache=shared")
	if err != nil {
		return nil, err
	}
	db, err := bun.NewDB(sqldb, sqlitedialect.New()), nil

	// enable foreign key constraint
	db.Exec("PRAGMA foreign_keys = ON")

	// db execute log
	db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))
	return db, err
}

// init postgres db for production perpose
// TODO(wj, important, low): fill this function
// TODO(wj, low): migration
func InitPGDB() (*bun.DB, error) {
	return nil, fmt.Errorf("not implemented")
}

// init database models
// should be database agnostic
func InitModels(db *bun.DB) error {
	ctx := context.Background()

	db.RegisterModel((*PostTag)(nil))

	_, err := db.NewCreateTable().Model((*User)(nil)).Exec(ctx)
	if err != nil {
		return err
	}

	_, err = db.NewCreateTable().Model((*Post)(nil)).Exec(ctx)
	if err != nil {
		return err
	}

	_, err = db.NewCreateTable().Model((*Comment)(nil)).Exec(ctx)
	if err != nil {
		return err
	}

	_, err = db.NewCreateTable().Model((*Attachment)(nil)).Exec(ctx)
	if err != nil {
		return err
	}

	_, err = db.NewCreateTable().Model((*Follow)(nil)).Exec(ctx)
	if err != nil {
		return err
	}

	// new message model
	_, err = db.NewCreateTable().Model((*Message)(nil)).Exec(ctx)
	if err != nil {
		return err
	}

	// new bookmark model
	_, err = db.NewCreateTable().Model((*Bookmark)(nil)).Exec(ctx)
	if err != nil {
		return err
	}

	// new bookmark_list model
	_, err = db.NewCreateTable().Model((*BookmarkList)(nil)).Exec(ctx)
	if err != nil {
		return err
	}


	// new block model
	_, err = db.NewCreateTable().Model((*Block)(nil)).Exec(ctx)
	if err != nil {
		return err
	}

	// new tag model
	_, err = db.NewCreateTable().Model((*Tag)(nil)).Exec(ctx)
	if err != nil {
		return err
	}

	// new post_tag model
	_, err = db.NewCreateTable().Model((*PostTag)(nil)).Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (*Post) BeforeCreateTable(ctx context.Context, query *bun.CreateTableQuery) error {
	query.ForeignKey(`("post_userid") REFERENCES "user" ("user_id") ON DELETE CASCADE`)
	return nil
}

func (*Comment) BeforeCreateTable(ctx context.Context, query *bun.CreateTableQuery) error {
	query.ForeignKey(`("comment_userid") REFERENCES "user" ("user_id") ON DELETE CASCADE`)
	query.ForeignKey(`("comment_postid") REFERENCES "post" ("post_id") ON DELETE CASCADE`)
	return nil
}

// why primary key not working?
// cos sqlite needs PRAGMA foreign_keys = ON;
// TODO(wj, mid): add index to speed up query
