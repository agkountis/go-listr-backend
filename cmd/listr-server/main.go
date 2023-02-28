package main

import (
	"github.com/agkountis/go-listr-backend/internal/db"
	"github.com/agkountis/go-listr-backend/internal/endpoints"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.Use(func(ctx *gin.Context) {
		ctx.Set("db", db.OpenConnection())
	}).
		POST("/lists", endpoints.CreateList).
		GET("/lists", endpoints.GetLists)
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
