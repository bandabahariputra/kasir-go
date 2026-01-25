package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Product struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
	Stock int    `json:"stock"`
}

var products = []Product{
	{
		ID:    1,
		Name:  "Indomie",
		Price: 3500,
		Stock: 10,
	},
	{
		ID:    2,
		Name:  "C 1000",
		Price: 5000,
		Stock: 15,
	},
}

type Category struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

var categories = []Category{
	{
		ID:          1,
		Name:        "Makanan",
		Description: "Ini adalah barang makanan",
	},
	{
		ID:          2,
		Name:        "Minuman",
		Description: "Ini adalah barang minuman",
	},
}

func main() {
	// GET http://localhost:8080/api/products/{id}
	// PUT http://localhost:8080/api/products/{id}
	// DELETE http://localhost:8080/api/products/{id}
	http.HandleFunc("/api/products/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			getProductById(w, r)
		case "PUT":
			updateProductById(w, r)
		case "DELETE":
			deleteProductById(w, r)
		}
	})

	// GET http://localhost:8080/api/product
	// POST http://localhost:8080/api/product
	http.HandleFunc("/api/products", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			getProducts(w, r)
		case "POST":
			createProduct(w, r)
		}
	})

	// GET http://localhost:8080/api/categories/{id}
	// PUT http://localhost:8080/api/categories/{id}
	// DELETE http://localhost:8080/api/categories/{id}
	http.HandleFunc("/api/categories/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			getCategoryById(w, r)
		case "PUT":
			updateCategoryById(w, r)
		case "DELETE":
			deleteCategoryById(w, r)
		}
	})

	// GET http://localhost:8080/api/categories
	// POST http://localhost:8080/api/categories
	http.HandleFunc("/api/categories", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			getCategories(w, r)
		case "POST":
			createCategory(w, r)
		}
	})

	// GET http://localhost:8080/health
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status": "up",
		})
	})

	fmt.Println("Server starting on port 8080...")

	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

// products
// get products
func getProducts(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

// create products
func createProduct(w http.ResponseWriter, r *http.Request) {
	var newProduct Product
	err := json.NewDecoder(r.Body).Decode(&newProduct)

	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if newProduct.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	if newProduct.Price <= 0 || newProduct.Stock <= 0 {
		http.Error(w, "Price and Stock are required", http.StatusBadRequest)
		return
	}

	newProduct.ID = len(products) + 1
	products = append(products, newProduct)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newProduct)
}

// get product by id
func getProductById(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/products/")

	id, err := strconv.Atoi(idStr)

	if err != nil {
		http.Error(w, "Invalid product id", http.StatusBadRequest)
		return
	}

	for _, p := range products {
		if p.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(p)
			return
		}
	}

	http.Error(w, "Product not found", http.StatusNotFound)
}

// update product by id
func updateProductById(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/products/")

	id, err := strconv.Atoi(idStr)

	if err != nil {
		http.Error(w, "Invalid product id", http.StatusBadRequest)
		return
	}

	var newProduct Product
	err = json.NewDecoder(r.Body).Decode(&newProduct)

	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if newProduct.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	if newProduct.Price <= 0 || newProduct.Stock <= 0 {
		http.Error(w, "Price and Stock are required", http.StatusBadRequest)
		return
	}

	for i := range products {
		if products[i].ID == id {
			newProduct.ID = id
			products[i] = newProduct

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(newProduct)
			return
		}
	}

	http.Error(w, "Product not found", http.StatusNotFound)
}

// delete product by id
func deleteProductById(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/products/")

	id, err := strconv.Atoi(idStr)

	if err != nil {
		http.Error(w, "Invalid product id", http.StatusBadRequest)
		return
	}

	for i, p := range products {
		if p.ID == id {
			products = append(products[:i], products[i+1:]...)

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"message": "Product deleted",
			})
			return
		}
	}

	http.Error(w, "Product not found", http.StatusNotFound)
}

// categories
// get categories
func getCategories(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}

// create category
func createCategory(w http.ResponseWriter, r *http.Request) {
	var newCategory Category
	err := json.NewDecoder(r.Body).Decode(&newCategory)

	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if newCategory.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	if newCategory.Description == "" {
		http.Error(w, "Description is required", http.StatusBadRequest)
		return
	}

	newCategory.ID = len(categories) + 1
	categories = append(categories, newCategory)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newCategory)
}

// get category by id
func getCategoryById(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")

	id, err := strconv.Atoi(idStr)

	if err != nil {
		http.Error(w, "Invalid category id", http.StatusBadRequest)
		return
	}

	for _, c := range categories {
		if c.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(c)
			return
		}
	}

	http.Error(w, "Category not found", http.StatusNotFound)
}

// update category by id
func updateCategoryById(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")

	id, err := strconv.Atoi(idStr)

	if err != nil {
		http.Error(w, "Invalid category id", http.StatusBadRequest)
		return
	}

	var newCategory Category
	err = json.NewDecoder(r.Body).Decode(&newCategory)

	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if newCategory.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	if newCategory.Description == "" {
		http.Error(w, "Description is required", http.StatusBadRequest)
		return
	}

	for i := range categories {
		if categories[i].ID == id {
			newCategory.ID = id
			categories[i] = newCategory

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(newCategory)
			return
		}
	}

	http.Error(w, "Category not found", http.StatusNotFound)
}

// delete category by id
func deleteCategoryById(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")

	id, err := strconv.Atoi(idStr)

	if err != nil {
		http.Error(w, "Invalid category id", http.StatusBadRequest)
		return
	}

	for i, c := range categories {
		if c.ID == id {
			categories = append(categories[:i], categories[i+1:]...)

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"message": "Category deleted",
			})
			return
		}
	}

	http.Error(w, "Category not found", http.StatusNotFound)
}
