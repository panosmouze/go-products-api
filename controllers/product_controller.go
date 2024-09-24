package controllers

import (
    "encoding/json"
    "strconv"
    "errors"
    "go-products-api/database"
    "go-products-api/models"
    "net/http"
    "github.com/gorilla/mux"
)

func validateProduct(product models.Product) error {
    if product.Name == "" {
        return errors.New("product name is required")
    }
    return nil
}

// POST /product
func CreateProduct(w http.ResponseWriter, r *http.Request) {
    var product models.Product

    if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    if err := validateProduct(product); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    if err := database.DB.Create(&product).Error; err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(product)
}

// GET /products
func GetProducts(w http.ResponseWriter, r *http.Request) {
    var page int = 1
    var limit int = 5
    var offset int
    var products []models.Product
	var totalProducts int64

    pageStr := r.URL.Query().Get("page")
    limitStr := r.URL.Query().Get("limit")

    if pageStr != "" {
        if parsedPage, err := strconv.ParseUint(pageStr, 10, 32); err == nil {
            page = int(parsedPage)
        }
    }

    if limitStr != "" {
        if parsedLimit, err := strconv.ParseUint(limitStr, 10, 32); err == nil {
            limit = int(parsedLimit)
        }
    }

	offset = (page - 1) * limit

	database.DB.Model(&models.Product{}).Count(&totalProducts)
	maxPage := int((totalProducts + int64(limit) - 1) / int64(limit))
	database.DB.Limit(limit).Offset(offset).Find(&products)

    w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"products": products,
		"page":     page,
		"maxPage":  maxPage,
		"total":    totalProducts,
	})
}

// GET /product/{id}
func GetProductByID(w http.ResponseWriter, r *http.Request) {
    var product models.Product

    id := mux.Vars(r)["id"]

    if err := database.DB.First(&product, "id = ?", id).Error; err != nil {
        http.Error(w, "Product not found", http.StatusNotFound)
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(product)
}

// PUT /product/{id}
func UpdateProduct(w http.ResponseWriter, r *http.Request) {
    var product models.Product

    id := mux.Vars(r)["id"]

    if err := database.DB.First(&product, "id = ?", id).Error; err != nil {
        http.Error(w, "Product not found", http.StatusNotFound)
        return
    }

    if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    if err := validateProduct(product); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    database.DB.Save(&product)
    w.WriteHeader(http.StatusNoContent)
    json.NewEncoder(w).Encode(product)
}

// DELETE /product/{id}
func DeleteProduct(w http.ResponseWriter, r *http.Request) {
    var product models.Product

    id := mux.Vars(r)["id"]

    if err := database.DB.First(&product, "id = ?", id).Error; err != nil {
        http.Error(w, "Product not found", http.StatusNotFound)
        return
    }

    database.DB.Delete(&product)
    w.WriteHeader(http.StatusNoContent)
}
