package main

import (
	"github.com/agkountis/go-listr-backend/internal/db"
	"github.com/agkountis/go-listr-backend/internal/model"
)

func main() {
	db := db.OpenConnection()
	db.AutoMigrate(&model.List{}, &model.Item{})
}
