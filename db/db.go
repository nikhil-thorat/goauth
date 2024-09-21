package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/nikhil-thorat/goauth/configs"
)

func NewPostgresStorage() (*sql.DB, error) {

	db, err := sql.Open("postgres", configs.Envs.DBConnectionString)
	if err != nil {
		return nil, fmt.Errorf("FAILED TO OPEN DATABASE : %v", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("FAILED TO PING THE DATABASE : %v", err)
	}

	log.Println("DATABASE CONNECTION SUCCESSFUL")
	return db, nil
}
