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

	return nil
}
