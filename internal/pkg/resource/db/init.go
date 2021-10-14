/*
 * @Author: Adrian Faisal
 * @Date: 14/10/21 13.44
 */

package db

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/keleeeep/test/internal/pkg/model"
	_ "github.com/mattn/go-sqlite3"
)

type Persistent interface {
	CreateUser(ctx context.Context, user *model.User) (*model.User, error)
	FindUser(ctx context.Context, data, column string) (*model.User, error)
}

type persistent struct {
	conn *sql.DB
}

func NewPersistent(datasource string) (Persistent, error) {
	db, err := sql.Open("sqlite3", datasource)
	if err != nil {
		return nil, fmt.Errorf("open database connection failed: %v", err)
	}

	statement, err := db.Prepare("CREATE TABLE IF NOT EXISTS users(id INTEGER PRIMARY KEY AUTOINCREMENT, name varchar, phone varchar, password varchar, role varchar, timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP)")
	if err != nil {
		return nil, fmt.Errorf("create table failed: %v", err)
	}
	statement.Exec()

	// ping a database connection is recommended to verify the connection still alive
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("ping database failed: %v", err)
	}

	return &persistent{conn: db}, nil
}
