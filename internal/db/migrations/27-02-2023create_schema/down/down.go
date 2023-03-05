package main

import (
	"github.com/agkountis/go-listr-backend/internal/db"
	"github.com/agkountis/go-listr-backend/internal/model"
)

func main() {
	db := db.OpenConnection()

	db.Migrator().DropTable(&model.List{}, &model.ListItem{})
}
