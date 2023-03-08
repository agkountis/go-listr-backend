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

	migrator := db.Migrator()

	if migrator.HasTable(&database.List{}) {
		migrator.DropTable(&database.List{})
	}

	if migrator.HasTable(&database.ListItem{}) {
		migrator.DropTable(&database.ListItem{})
	}

	db.AutoMigrate(&database.List{}, &database.ListItem{})
}
