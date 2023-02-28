package db

import (
	"context"
	"database/sql"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
	"github.com/uptrace/bun/extra/bundebug"
)

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

func InitModels(db *bun.DB) error {
	ctx := context.Background()
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

func (*Attachment) BeforeCreateTable(ctx context.Context, query *bun.CreateTableQuery) error {
	query.ForeignKey(`("attachment_postid") REFERENCES "post" ("post_id") ON DELETE CASCADE`)
	return nil
}

// why primary key not working?
// cos sqlite needs PRAGMA foreign_keys = ON;
