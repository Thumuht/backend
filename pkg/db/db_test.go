package db

import (
	"context"
	"testing"

	"github.com/uptrace/bun"
)

var db *bun.DB

func TestMain(m *testing.M) {
	db, _ = InitSQLiteDB()
	InitModels(db)
	m.Run()
}

func TestDB(t *testing.T) {
	ctx := context.Background()
	user := &User{
		LoginName: "thumuht",
		Nickname:  "THUMUHT",
		Email:     "a@a.com",
		About:     "the man without quality",
	}
	_, err := db.NewInsert().Model(user).Exec(ctx)
	if err != nil {
		t.Errorf("insert user failed.")
	}

	dbUser := new(User)
	err = db.NewSelect().Model(dbUser).Where("login_name = ?", "thumuht").Scan(ctx)
	if err != nil {
		t.Errorf("query user failed.")
	}

	if dbUser.Nickname != user.Nickname || dbUser.Email != user.Email ||
		dbUser.About != user.About {
		t.Errorf("query field not identical to original one")
	}

}
