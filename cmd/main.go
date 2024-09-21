package main

import (
	"log"

	"github.com/nikhil-thorat/goauth/cmd/server"
	"github.com/nikhil-thorat/goauth/configs"
	"github.com/nikhil-thorat/goauth/db"
)

func main() {

	db, err := db.NewPostgresStorage()
	if err != nil {
		log.Fatal("DATABASE ERROR : ", err)
	}
	defer db.Close()

	server := server.NewServer(configs.Envs.Port, db)

	err = server.Run()
	if err != nil {
		log.Fatal("ERROR STARTING SERVER : ", err)
	}

}
