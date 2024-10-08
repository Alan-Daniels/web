package database

import (
	"github.com/Alan-Daniels/web/internal/config"
	"github.com/surrealdb/surrealdb.go"
)

type DB struct {
	db *surrealdb.DB
}

type Map map[string]interface{}

func Init(Config *config.Config) (*DB, error) {
	db, err := surrealdb.New(Config.Database.Uri)
	if err != nil {
		return nil, err
	}

	authData := map[string]interface{}{
		"NS":   Config.Database.Namespace,
		"user": Config.Database.Username,
		"pass": Config.Database.Password,
	}
	if _, err = db.Signin(authData); err != nil {
		return nil, err
	}
	Config.Database.Username = ""
	Config.Database.Password = ""

	if _, err = db.Use(Config.Database.Namespace, Config.Database.Name); err != nil {
		return nil, err
	}
	return &DB{db: db}, nil
}

func (db *DB) RootBranches() (interface{}, error) {
	return db.db.Query("SELECT * FROM Branch WHERE parent IS NULL", Map{})
}
func (db *DB) Branches(parent string) (interface{}, error) {
	return db.db.Query("SELECT * FROM Branch WHERE (parent=$parent)", Map{
		"parent": parent,
	})
}

func (db *DB) Query(sql string, vars interface{}) (interface{}, error) {
	return db.db.Query(sql, vars)
}
