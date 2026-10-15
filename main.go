package main

import (
	"final-project/pkg/db"
	"final-project/pkg/server"
	"log"
)

func main() {

	err := db.Init("scheduler.db")
	if err != nil {
		log.Printf("Ошибка %v", err)
		return
	}

	defer db.Close()

	srv := server.CreateServer()
	if err := srv.Server.ListenAndServe(); err != nil {
		log.Printf("Ошибка %v", err)
	}

}
