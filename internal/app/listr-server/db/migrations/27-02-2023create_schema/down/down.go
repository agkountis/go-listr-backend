package main

import (
	"github.com/agkountis/go-listr-backend/internal/app/listr-server/db"
	"github.com/agkountis/go-listr-backend/internal/app/listr-server/model"
)

func main() {
	db := db.OpenConnection()

	db.Migrator().DropTable(&model.List{}, &model.ListItem{})
}
