package main

import (
    "go-products-api/config"
    "go-products-api/database"
    "net/http"
)

func main() {
    database.Connect()

    router := config.SetupRoutes()

    http.ListenAndServe(":8080", router)
}