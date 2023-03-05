package main

import (
	"github.com/agkountis/go-listr-backend/internal/db"
	"github.com/agkountis/go-listr-backend/internal/endpoints"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Self signed certs are a paint to work with. //TODO
// var selfSignedCertsPath string = os.Getenv("DOMAIN_SELF_SIGNED_CERTS_PATH")
// var certFilePath string = fmt.Sprintf("%v/domain.crt", selfSignedCertsPath)
// var keyFilePath string = fmt.Sprintf("%v/domain.key", selfSignedCertsPath)

func main() {
	var DB *gorm.DB = db.OpenConnection()

	r := gin.Default()

	// Middleware that provides a ref to the DB connection pool for all endpoints
	r.Use(func(ctx *gin.Context) {
		ctx.Set("db", DB)
	})

	v1 := r.Group("/v1")

	v1.POST("/lists", endpoints.CreateList)
	v1.POST("/lists/:id", endpoints.CreateListItem)

	v1.GET("/lists", endpoints.GetLists)
	v1.GET("/lists/:id", endpoints.GetListItems)

	v1.DELETE("lists", endpoints.DeleteList)
	v1.DELETE("lists/:id", endpoints.DeleteListItem)

	// Self signed certs are a pain to work with. //TODO
	// r.RunTLS("0.0.0.0:8080", certFilePath, keyFilePath)
	r.Run("0.0.0.0:8080") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
