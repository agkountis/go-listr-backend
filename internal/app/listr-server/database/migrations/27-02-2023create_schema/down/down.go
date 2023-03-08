package main

import (
	"log"

	"github.com/agkountis/go-listr-backend/internal/app/listr-server/database"
)

func main() {
	context, err := database.OpenConnection()
	if err != nil {
		log.Fatalf("Failed to connect to database. Error: %v", err.Error())
	}
	db := context
	db.Migrator().DropTable(&database.List{}, &database.ListItem{})
}
