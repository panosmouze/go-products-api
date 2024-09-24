package database

import (
    "log"
    "os"
    "go-products-api/models"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
    var err error

    appEnv := os.Getenv("APP_ENV")
    if appEnv == "" {
        appEnv = "development"
    }

    switch appEnv {
    case "development":
        DB, err = gorm.Open(sqlite.Open("development.db"), &gorm.Config{})
        if err != nil {
            log.Fatal("Failed to connect to SQLite in development:", err)
        }
        log.Println("Using SQLite in development mode")
        
    case "production":
        DB, err = gorm.Open(sqlite.Open("production.db"), &gorm.Config{})
        if err != nil {
            log.Fatal("Failed to connect to SQLite in production:", err)
        }
        log.Println("Using SQLite in production mode")
        
    default:
        log.Fatalf("Unknown APP_ENV: %s", appEnv)
    }

    // Migrate the database schema
    err = DB.AutoMigrate(&models.Product{})
    if err != nil {
        log.Fatal("Failed to migrate the database schema:", err)
    }
}
