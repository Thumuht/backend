package db

import (
	"context"
	"database/sql"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
)

func InitSQLiteDB() (*bun.DB, error) {
	sqldb, err := sql.Open(sqliteshim.ShimName, "file::memory:?cache=shared")
	if err != nil {
		return nil, err
	}
	return bun.NewDB(sqldb, sqlitedialect.New()), nil
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
