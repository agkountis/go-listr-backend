package database

import (
	"fmt"
	"log"
	"os"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var connectionStringTemplate string = "postgresql://%v:%v@ring-raven-4280.6zw.cockroachlabs.cloud:26257/listr?sslmode=verify-full"
var sqlUserName string = os.Getenv("COCKROACHDB_LISTR_UN")
var sqlUserPasswd string = os.Getenv("COCKROACHDB_LISTR_PWD")

func OpenConnection() (*gorm.DB, error) {
	dsn := fmt.Sprintf(connectionStringTemplate, sqlUserName, sqlUserPasswd)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("failed to connect database. Err: %v", err)
		return nil, err
	}

	log.Println("DB connection successful. Connection pool initialized...")
	return db, nil
}

type List struct {
	ID   uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	Name string    `gorm:"not null"`
}

type ListItem struct {
	ID     uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	Data   string    `gorm:"not null"`
	ListID uuid.UUID `gorm:"not null"`
	List   List      `json:"-"`
}
