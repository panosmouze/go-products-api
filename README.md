# Go Products API

A simple RESTful API for managing products using Go.

## Clone the Repository

To download the code from the GitHub repository, run the following command:

```bash
git clone https://github.com/panosmouze/go-products-api.git
cd go-products-api
go mod tidy
go build -o myapp
./myapp
```

### Build and Run
To build and run the app, run the following command:

```bash
go mod tidy
go build -o myapp
./myapp
```

### Run unit tests
To run the unit tests, run the following command:

```bash
go test -v
```

## Running with Docker
Pull the Image from Docker Hub

```bash
docker pull panosmouze/go-products-api 
docker run -d -p 8080:8080 --name go-products-api-container panosmouze/go-products-api
```

The server will be accessible at http://localhost:8080.

## Manual testing with Postman Collection

To quickly test the API, you can use the provided Postman collection located in this repo under /docs.

## API Endpoints Overview

1. **Get Product by ID**
   - **Endpoint**: `GET /product/{id}`
   - **Description**: Retrieves a specific product using its unique ID.
   - **Example Request**:
     ```
     GET http://localhost:8080/product/1
     ```

2. **Get Products (Paginated)**
   - **Endpoint**: `GET /products`
   - **Description**: Retrieves a list of products with pagination. You can specify the page number and limit the number of products returned.
   - **Query Parameters**:
     - `page`: The page number to retrieve (default: 1).
     - `limit`: The maximum number of products to return per page (default: 5).
   - **Example Request**:
     ```
     GET http://localhost:8080/products?page=1&limit=5
     ```

3. **Create a New Product**
   - **Endpoint**: `POST /product`
   - **Description**: Creates a new product. You must send a JSON body containing the product data.
   - **Request Body** (example):
     ```json
     {
       "name": "product name"
     }
     ```
   - **Example Request**:
     ```
     POST http://localhost:8080/product
     ```

4. **Update an Existing Product**
   - **Endpoint**: `PUT /product/{id}`
   - **Description**: Updates the details of an existing product specified by its unique ID. You must send a JSON body with the updated product data.
   - **Request Body** (example):
     ```json
     {
       "name": "new product name"
     }
     ```
   - **Example Request**:
     ```
     PUT http://localhost:8080/product/1
     ```

5. **Delete a Product**
   - **Endpoint**: `DELETE /product/{id}`
   - **Description**: Deletes a specific product identified by its unique ID.
   - **Example Request**:
     ```
     DELETE http://localhost:8080/product/1
     ```
