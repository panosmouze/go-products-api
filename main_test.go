package main

import (
    "os"
    "log"
    "bytes"
    "net/http"
    "net/http/httptest"
    "testing"
    "go-products-api/database"
    "go-products-api/config"
    "encoding/json"
    "fmt"
    "github.com/gorilla/mux"
)

type ProductResponse struct {
    ID   uint   `json:"id"`
    Name string `json:"name"`
}

type ProductsPageResponse struct {
    Products    []ProductResponse   `json:"products"`
    Page        int                 `json:"page"`
    MaxPage     int                 `json:"maxPage"`
    Total       int                 `json:"total"`
}

var router *mux.Router

func TestMain(m *testing.M) {
    if err := os.Remove("development.db"); err != nil && !os.IsNotExist(err) {
        log.Fatalf("Failed to remove database: %v", err)
    }

    database.Connect()
    router = config.SetupRoutes()

    exitCode := m.Run()

    os.Exit(exitCode)
}

func Test1_CreateProducts(t *testing.T) {
    t.Run("Create 100 valid products", func(t *testing.T) {
        for i := 0; i < 100; i++ {
            jsonStr := []byte(`{"name":"tempname"}`)
            rr := RequestCreateProduct(jsonStr, t)
            ValidateCreateProduct(rr, http.StatusCreated, t)
        }
    })

    t.Run("Create product with empty name", func(t *testing.T) {
        jsonStr := []byte(`{"name":""}`)
        rr := RequestCreateProduct(jsonStr, t)
        ValidateCreateProduct(rr, http.StatusBadRequest, t)
    })

    t.Run("Create product with missing name", func(t *testing.T) {
        jsonStr := []byte(`{}`)
        rr := RequestCreateProduct(jsonStr, t)
        ValidateCreateProduct(rr, http.StatusBadRequest, t)
    })
}

func Test2_GetProduct(t *testing.T) {
    t.Run("Get product with id 1", func(t *testing.T) {
        rr := RequestGetProduct(1, t)
        ValidateGetProduct(1, "tempname", rr, http.StatusOK, t)
    })
    t.Run("Get product with id 101", func(t *testing.T) {
        rr := RequestGetProduct(101, t)
        ValidateGetProduct(101, "", rr, http.StatusNotFound, t)
    })
}

func Test3_UpdateProduct(t *testing.T) {
    t.Run("Update product with id 1.1", func(t *testing.T) {
        jsonStr := []byte(`{"name":"updatedname"}`)
        rr := RequestUpdateProduct(1, jsonStr, t)
        ValidateUpdateProduct(rr, http.StatusNoContent, t)

        rr = RequestGetProduct(1, t)
        ValidateGetProduct(1, "updatedname", rr, http.StatusOK, t)
    })

    t.Run("Update product with id 1.2", func(t *testing.T) {
        jsonStr := []byte(`{}`)
        rr := RequestUpdateProduct(1, jsonStr, t)
        ValidateUpdateProduct(rr, http.StatusNoContent, t)
    })

    t.Run("Update product with id 101", func(t *testing.T) {
        jsonStr := []byte(`{"name":"newname"}`)
        rr := RequestUpdateProduct(101, jsonStr, t)
        ValidateUpdateProduct(rr, http.StatusNotFound, t)
    })

    t.Run("Update product with id 1 with empty name", func(t *testing.T) {
        jsonStr := []byte(`{"name":""}`)
        rr := RequestUpdateProduct(1, jsonStr, t)
        ValidateUpdateProduct(rr, http.StatusBadRequest, t)
    })
}

func Test4_DeleteProduct(t *testing.T) {
    t.Run("Delete product with id 1", func(t *testing.T) {
        rr := RequestDeleteProduct(1, t)
        ValidateDeleteProduct(rr, http.StatusNoContent, t)

        rr = RequestGetProduct(1, t)
        ValidateGetProduct(1, "", rr, http.StatusNotFound, t)
    })

    t.Run("Delete product with id 101", func(t *testing.T) {
        rr := RequestDeleteProduct(101, t)
        ValidateDeleteProduct(rr, http.StatusNotFound, t)
    })
}

func Test5_GetProductsPaginate(t *testing.T) {
    t.Run("Get products with pagination page 1 limit 10", func(t *testing.T) {
        rr := RequestGetProductsPaginated(1, 10, t)
        ValidateGetProductsPaginated(10, 99, rr, http.StatusOK, t)
    })

    t.Run("Get products with pagination page 10 limit 10", func(t *testing.T) {
        rr := RequestGetProductsPaginated(10, 10, t)
        ValidateGetProductsPaginated(9, 99, rr, http.StatusOK, t)
    })

    t.Run("Get products with pagination page 100 limit 100", func(t *testing.T) {
        rr := RequestGetProductsPaginated(100, 100, t)
        ValidateGetProductsPaginated(0, 99, rr, http.StatusOK, t)
    })
}

func RequestCreateProduct(jsonStr []byte, t *testing.T) *httptest.ResponseRecorder {
    req, err := http.NewRequest(http.MethodPost, "/product", bytes.NewBuffer(jsonStr))
    if err != nil {
        t.Fatalf("Failed to create request: %v", err)
    }
    req.Header.Set("Content-Type", "application/json")

    rr := httptest.NewRecorder()
    router.ServeHTTP(rr, req)

    return rr
}

func ValidateCreateProduct(rr *httptest.ResponseRecorder, expectedStatusCode int, t *testing.T) {
    if status := rr.Code; status != expectedStatusCode {
        t.Errorf("Handler returned wrong status code: got %v want %v", status, expectedStatusCode)
    }

    if expectedStatusCode == http.StatusCreated {
        var product ProductResponse
        if err := json.NewDecoder(rr.Body).Decode(&product); err != nil {
            t.Fatalf("Failed to decode response: %v", err)
        }

        if product.ID <= 0 {
            t.Errorf("Expected a valid product ID, but got %d", product.ID)
        }

        if product.Name != "tempname" {
            t.Errorf("Expected product name to be 'tempname', but got '%s'", product.Name)
        }
    }
}

func RequestGetProduct(id uint, t *testing.T) *httptest.ResponseRecorder {
    url := fmt.Sprintf("/product/%d", id)
    req, err := http.NewRequest(http.MethodGet, url, nil)
    if err != nil {
        t.Fatalf("Failed to create request: %v", err)
    }

    rr := httptest.NewRecorder()
    router.ServeHTTP(rr, req)

    return rr
}

func ValidateGetProduct(expectedID uint, expectedName string, rr *httptest.ResponseRecorder, expectedStatusCode int, t *testing.T) {
    if status := rr.Code; status != expectedStatusCode {
        t.Errorf("Handler returned wrong status code: got %v want %v", status, expectedStatusCode)
    }

    if expectedStatusCode == http.StatusOK {
        var product ProductResponse
        if err := json.NewDecoder(rr.Body).Decode(&product); err != nil {
            t.Fatalf("Failed to decode response: %v", err)
        }

        if product.ID != expectedID {
            t.Errorf("Expected product ID %d but got %d", expectedID, product.ID)
        }

        if product.Name != expectedName {
            t.Errorf("Expected product name to be '%s', but got '%s'", expectedName, product.Name)
        }
    }
}

func RequestUpdateProduct(id uint, jsonStr []byte, t *testing.T) *httptest.ResponseRecorder {
    url := fmt.Sprintf("/product/%d", id)
    req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(jsonStr))
    if err != nil {
        t.Fatalf("Failed to create request: %v", err)
    }
    req.Header.Set("Content-Type", "application/json")

    rr := httptest.NewRecorder()
    router.ServeHTTP(rr, req)

    return rr
}

func ValidateUpdateProduct(rr *httptest.ResponseRecorder, expectedStatusCode int, t *testing.T) {
    if status := rr.Code; status != expectedStatusCode {
        t.Errorf("Handler returned wrong status code: got %v want %v", status, expectedStatusCode)
    }
}

func RequestDeleteProduct(id uint, t *testing.T) *httptest.ResponseRecorder {
    url := fmt.Sprintf("/product/%d", id)
    req, err := http.NewRequest(http.MethodDelete, url, nil)
    if err != nil {
        t.Fatalf("Failed to create request: %v", err)
    }

    rr := httptest.NewRecorder()
    router.ServeHTTP(rr, req)

    return rr
}

func ValidateDeleteProduct(rr *httptest.ResponseRecorder, expectedStatusCode int, t *testing.T) {
    if status := rr.Code; status != expectedStatusCode {
        t.Errorf("Handler returned wrong status code: got %v want %v", status, expectedStatusCode)
    }
}

func RequestGetProductsPaginated(page, limit int, t *testing.T) *httptest.ResponseRecorder {
    url := fmt.Sprintf("/products?page=%d&limit=%d", page, limit)
    req, err := http.NewRequest(http.MethodGet, url, nil)
    if err != nil {
        t.Fatalf("Failed to create request: %v", err)
    }

    rr := httptest.NewRecorder()
    router.ServeHTTP(rr, req)

    return rr
}

func ValidateGetProductsPaginated(expectedProds, totalProds int, rr *httptest.ResponseRecorder, expectedStatusCode int, t *testing.T) {
    if status := rr.Code; status != expectedStatusCode {
        t.Errorf("Handler returned wrong status code: got %v want %v", status, expectedStatusCode)
    }

    if expectedStatusCode == http.StatusOK {
        var page ProductsPageResponse
        if err := json.NewDecoder(rr.Body).Decode(&page); err != nil {
            t.Fatalf("Failed to decode response: %v", err)
        }

        if len(page.Products) != expectedProds {
            t.Errorf("Expected products %d but got %d", expectedProds, len(page.Products))
        }

        if page.Total != totalProds {
            t.Errorf("Expected products %d but got %d", totalProds, page.Total)
        }
    }
}
