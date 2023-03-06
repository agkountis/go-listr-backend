package main

import (
	"github.com/agkountis/go-listr-backend/internal/app/listr-server/db"
	"github.com/agkountis/go-listr-backend/internal/app/listr-server/model"
)

func main() {
	db := db.OpenConnection()
	migrator := db.Migrator()

	if migrator.HasTable(&model.List{}) {
		migrator.DropTable(&model.List{})
	}

	if migrator.HasTable(&model.ListItem{}) {
		migrator.DropTable(&model.ListItem{})
	}

	db.AutoMigrate(&model.List{}, &model.ListItem{})
}
