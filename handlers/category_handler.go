package handlers

import (
	"encoding/json"
	"kasir-go/models"
	"kasir-go/services"
	"net/http"
	"strconv"
	"strings"
)

type CategoryHandler struct {
	service *services.CategoryService
}

func NewCategoryHandler(service *services.CategoryService) *CategoryHandler {
	return &CategoryHandler{service: service}
}

func (h *CategoryHandler) HandleCategories(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetAll(w, r)
	case http.MethodPost:
		h.Create(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

// GET http://localhost:8080/api/categories
func (h *CategoryHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	categories, err := h.service.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}

// POST http://localhost:8080/api/categories
func (h *CategoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	var category models.Category

	err := json.NewDecoder(r.Body).Decode(&category)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if category.Name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}

	if category.Description == "" {
		http.Error(w, "description is required", http.StatusBadRequest)
		return
	}

	err = h.service.Create(&category)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(category)
}

func (h *CategoryHandler) HandleCategoryByID(w http.ResponseWriter, r *http.Request) {
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

// GET http://localhost:8080/api/categories/{id}
func (h *CategoryHandler) GetById(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid category id", http.StatusBadRequest)
		return
	}

	category, err := h.service.GetById(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(category)
}

// PUT http://localhost:8080/api/categories/{id}
func (h *CategoryHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid category id", http.StatusBadRequest)
		return
	}

	var category models.Category
	err = json.NewDecoder(r.Body).Decode(&category)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if category.Name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}

	if category.Description == "" {
		http.Error(w, "description is required", http.StatusBadRequest)
		return
	}

	category.ID = id
	err = h.service.Update(&category)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(category)
}

// DELETE http://localhost:8080/api/categories/{id}
func (h *CategoryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid category id", http.StatusBadRequest)
		return
	}

	err = h.service.Delete(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Category deleted",
	})
}
