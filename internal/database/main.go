package database

import (
	"github.com/Alan-Daniels/web/internal/config"
	"github.com/surrealdb/surrealdb.go"
	"github.com/surrealdb/surrealdb.go/pkg/models"
)

type DB struct {
	db *surrealdb.DB
}

func Init(Config *config.Config) (*DB, error) {
	db, err := surrealdb.New(Config.Database.Uri)
	if err != nil {
		return nil, err
	}

	authData := models.Auth{
		Namespace: Config.Database.Namespace,
		Username:  Config.Database.Username,
		Password:  Config.Database.Password,
	}
	if _, err = db.Signin(&authData); err != nil {
		return nil, err
	}
	Config.Database.Username = ""
	Config.Database.Password = ""

	if _, err = db.Use(Config.Database.Namespace, Config.Database.Name); err != nil {
		return nil, err
	}
	return &DB{db: db}, nil
}

func (db *DB) Branches(parent string) (RawResponse, error) {
	resps, err := db.Query("SELECT * FROM Group WHERE parent = $parent", Map{
		"parent": parent,
	})
	if err != nil {
		return nil, err
	}
	return resps[0], nil
}

func (db *DB) Pages(parent string) (RawResponse, error) {
	resps, err := db.Query("SELECT * FROM Page WHERE parent = $parent", Map{
		"parent": parent,
	})
	if err != nil {
		return nil, err
	}
	return resps[0], nil
}

func (db *DB) Insert(table string, item interface{}) (RawResponse, error) {
	resps, err := toRawResponses(db.db.Insert(table, item))
	if err != nil {
		return nil, err
	}
	return resps[0], nil
}

func (db *DB) Query(sql string, vars interface{}) ([]RawResponse, error) {
	return toRawResponses(db.db.Query(sql, vars))
}
