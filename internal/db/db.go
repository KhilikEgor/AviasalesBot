package db

import (
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "log"
)

var DB *gorm.DB

func Connect() {
    dsn := "host=localhost user=admin password=admin dbname=aviasales_bot_bd port=5432 sslmode=disable"
    var err error
    DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatalf("failed to connect to the database: %v", err)
    }

    log.Println("Connected to the database successfully!")
}
