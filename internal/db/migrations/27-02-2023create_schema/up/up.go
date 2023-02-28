package main

import (
	"github.com/agkountis/go-listr-backend/internal/db"
	"github.com/agkountis/go-listr-backend/internal/model"
)

func main() {
	db := db.OpenConnection()
	migrator := db.Migrator()

	if migrator.HasTable(&model.List{}) {
		migrator.DropTable(&model.List{})
	}

	if migrator.HasTable(&model.Item{}) {
		migrator.DropTable(&model.Item{})
	}

	db.AutoMigrate(&model.List{}, &model.Item{})
}
