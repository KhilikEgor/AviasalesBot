package db

import (
    "gorm.io/gorm"
    "gorm.io/driver/postgres"
    "log"
)

var DB *gorm.DB // Глобальная переменная для хранения соединения

func Connect() {
    var err error
    dsn := "host=psql user=admin password=admin dbname=aviasales_bot_bd sslmode=disable"
    DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Printf("Failed to connect to the database with DSN: %s", dsn)
        log.Fatalf("failed to connect to the database: %v", err)
    }
    log.Println("Connected to the database successfully!")
}
