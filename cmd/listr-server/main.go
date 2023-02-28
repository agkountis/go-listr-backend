package main

import (
	"github.com/agkountis/go-listr-backend/internal/db"
	"github.com/agkountis/go-listr-backend/internal/endpoints"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func main() {
	var DB *gorm.DB = db.OpenConnection()

	r := gin.Default()

	// Middleware that provide a ref to the DB connection pool for all endpoints
	r.Use(func(ctx *gin.Context) {
		ctx.Set("db", DB)
	})

	v1 := r.Group("/v1")

	v1.POST("/lists", endpoints.CreateList)
	v1.GET("/lists", endpoints.GetLists)
	v1.GET("/lists/:id", endpoints.GetListItems)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
