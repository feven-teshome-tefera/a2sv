package main

import (
	"context"
	"log"
	"task_manager/db"
	"task_manager/router"
)

func main() {
	client := db.Connect()
	defer func() {
		ctx := context.Background()
		if err := client.Disconnect(ctx); err != nil {
			log.Fatal(err)
		}
		log.Println("MongoDB disconnected")
	}()

	r := router.SetupRouter()
	r.Run()
}
