package handlers

import (
	"encoding/json"
	"kasir-go/services"
	"net/http"
	"time"
)

type ReportHandler struct {
	service *services.ReportService
}

func NewReportHandler(service *services.ReportService) *ReportHandler {
	return &ReportHandler{service: service}
}

// GET http://localhost:8080/api/report/today
func (h *ReportHandler) GetTodayReport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	report, err := h.service.GetTodayReport()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(report)
}

// GET http://localhost:8080/api/report/today?start_date=YYYY-MM-DD&end_date=YYYY-MM-DD
func (h *ReportHandler) GetReport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	query := r.URL.Query()

	var startDatePtr *time.Time
	if startStr := query.Get("start_date"); startStr != "" {
		loc, _ := time.LoadLocation("Asia/Jakarta")
		startDate, err := time.ParseInLocation("2006-01-02", startStr, loc)
		if err != nil {
			http.Error(w, "invalid start_date format, use YYYY-MM-DD", http.StatusBadRequest)
			return
		}
		startDatePtr = &startDate
	}

	var endDatePtr *time.Time
	if endStr := query.Get("end_date"); endStr != "" {
		loc, _ := time.LoadLocation("Asia/Jakarta")
		endDate, err := time.ParseInLocation("2006-01-02", endStr, loc)
		if err != nil {
			http.Error(w, "invalid end_date format, use YYYY-MM-DD", http.StatusBadRequest)
			return
		}
		endDatePtr = &endDate
	}

	report, err := h.service.GetReport(startDatePtr, endDatePtr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(report)
}
