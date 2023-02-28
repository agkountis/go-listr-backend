package db

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var connectionStringTemplate string = "postgresql://%v:%v@ring-raven-4280.6zw.cockroachlabs.cloud:26257/listr?sslmode=verify-full"
var sqlUserName string = os.Getenv("COCKROACHDB_LISTR_UN")
var sqlUserPasswd string = os.Getenv("COCKROACHDB_LISTR_PWD")

func OpenConnection() *gorm.DB {
	dsn := fmt.Sprintf(connectionStringTemplate, sqlUserName, sqlUserPasswd)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database", err)
	}

	log.Println("DB connection successful. Connection pool initialized...")
	return db
}
