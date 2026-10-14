package main

import (
	"final-project/pkg/db"
	"final-project/pkg/server"
	"fmt"
)

func main() {

	err := db.Init("scheduler.db")
	if err != nil {
		fmt.Println(err)
		return
	}

	defer db.Close()

	srv := server.CreateServer()
	if err := srv.Server.ListenAndServe(); err != nil {
		fmt.Println(err)
	}

}
