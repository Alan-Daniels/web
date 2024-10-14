package database

import (
	"github.com/Alan-Daniels/web/internal/config"
	"github.com/surrealdb/surrealdb.go"
)

type DB struct {
	db *surrealdb.DB
}

func Init(Config *config.Config) (*surrealdb.DB, error) {
	db, err := surrealdb.New(Config.Database.Uri)
	if err != nil {
		return nil, err
	}

	authData := surrealdb.Auth{
		Namespace: Config.Database.Namespace,
		Username:  Config.Database.Username,
		Password:  Config.Database.Password,
	}
	if _, err = db.SignIn(&authData); err != nil {
		return nil, err
	}
	Config.Database.Username = ""
	Config.Database.Password = ""

	if err = db.Use(Config.Database.Namespace, Config.Database.Name); err != nil {
		return nil, err
	}

	return db, nil
}
