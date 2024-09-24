package config

import (
    "github.com/gorilla/mux"
    "go-products-api/controllers"
)

func SetupRoutes() *mux.Router {
    router := mux.NewRouter()

    router.HandleFunc("/product", controllers.CreateProduct).Methods("POST")
    router.HandleFunc("/products", controllers.GetProducts).Methods("GET")
    router.HandleFunc("/product/{id:[0-9]+}", controllers.GetProductByID).Methods("GET")
    router.HandleFunc("/product/{id:[0-9]+}", controllers.UpdateProduct).Methods("PUT")
    router.HandleFunc("/product/{id:[0-9]+}", controllers.DeleteProduct).Methods("DELETE")

    return router
}
