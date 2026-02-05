package handlers

import (
	"encoding/json"
	"kasir-go/models"
	"kasir-go/services"
	"net/http"
	"strconv"
	"strings"
)

type ProductHandler struct {
	service *services.ProductService
}

func NewProductHandler(service *services.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

func (h *ProductHandler) HandleProducts(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetAll(w, r)
	case http.MethodPost:
		h.Create(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

// GET http://localhost:8080/api/products
func (h *ProductHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")

	products, err := h.service.GetAll(name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

// POST http://localhost:8080/api/products
func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	var product models.Product

	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if product.Name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}

	if product.Price <= 0 {
		http.Error(w, "price is required", http.StatusBadRequest)
		return
	}

	if product.Stock <= 0 {
		http.Error(w, "stock are required", http.StatusBadRequest)
		return
	}

	if product.CategoryID == 0 {
		http.Error(w, "category id is required", http.StatusBadRequest)
		return
	}

	err = h.service.Create(&product)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(product)
}

func (h *ProductHandler) HandleProductByID(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetById(w, r)
	case http.MethodPut:
		h.Update(w, r)
	case http.MethodDelete:
		h.Delete(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

// GET http://localhost:8080/api/products/{id}
func (h *ProductHandler) GetById(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/products/")

	id, err := strconv.Atoi(idStr)

	if err != nil {
		http.Error(w, "invalid product id", http.StatusBadRequest)
		return
	}

	product, err := h.service.GetById(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

// PUT http://localhost:8080/api/products/{id}
func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/products/")

	id, err := strconv.Atoi(idStr)

	if err != nil {
		http.Error(w, "invalid product id", http.StatusBadRequest)
		return
	}

	var product models.Product
	err = json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if product.Name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}

	if product.Price <= 0 {
		http.Error(w, "price is required", http.StatusBadRequest)
		return
	}

	if product.Stock <= 0 {
		http.Error(w, "stock are required", http.StatusBadRequest)
		return
	}

	if product.CategoryID == 0 {
		http.Error(w, "category id is required", http.StatusBadRequest)
		return
	}

	product.ID = id
	err = h.service.Update(&product)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

// DELETE http://localhost:8080/api/products/{id}
func (h *ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/products/")

	id, err := strconv.Atoi(idStr)

	if err != nil {
		http.Error(w, "invalid product id", http.StatusBadRequest)
		return
	}

	err = h.service.Delete(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Product deleted",
	})
}
