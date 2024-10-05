package database

import (
	"github.com/Alan-Daniels/web/internal/config"
	"github.com/surrealdb/surrealdb.go"
)

// todo: put in its own module
func Init(Config *config.Config) (*surrealdb.DB, error) {
	db, err := surrealdb.New(Config.Database.Uri)
	if err != nil {
		return nil, err
	}

	authData := map[string]interface{}{
		"NS": Config.Database.Namespace,
		"user":  Config.Database.Username,
		"pass":  Config.Database.Password,
	}
	if _, err = db.Signin(authData); err != nil {
		return nil, err
	}
	Config.Database.Username = ""
	Config.Database.Password = ""

	if _, err = db.Use(Config.Database.Namespace, Config.Database.Name); err != nil {
		return nil, err
	}
	return db, nil
}
