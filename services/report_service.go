package services

import (
	"kasir-go/models"
	"kasir-go/repositories"
	"time"
)

type ReportService struct {
	repo *repositories.TransactionRepository
}

func NewReportService(repo *repositories.TransactionRepository) *ReportService {
	return &ReportService{repo: repo}
}

func (s *ReportService) GetTodayReport() (*models.TodayReport, error) {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	now := time.Now().In(loc)

	start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
	end := start.Add(24 * time.Hour)

	totalRevenue, totalTransaction, err := s.repo.GetSummaryByPeriod(start, end)
	if err != nil {
		return nil, err
	}

	productName, productQuantitySold, err := s.repo.GetBestSellingProductByPeriod(start, end)
	if err != nil {
		return nil, err
	}

	var bestProduct *models.BestSellingProduct
	if productName != "" || productQuantitySold > 0 {
		bestProduct = &models.BestSellingProduct{
			Name:         productName,
			QuantitySold: productQuantitySold,
		}
	}

	return &models.TodayReport{
		TotalRevenue:       totalRevenue,
		TotalTransaction:   totalTransaction,
		BestSellingProduct: bestProduct,
	}, nil
}

func (s *ReportService) GetReport(startDate, endDate *time.Time) (*models.TodayReport, error) {
	loc, _ := time.LoadLocation("Asia/Jakarta")

	var start time.Time
	if startDate != nil {
		start = startDate.In(loc)
	} else {
		start = time.Date(2026, time.January, 1, 0, 0, 0, 0, loc)
	}

	var end time.Time
	if endDate != nil {
		end = endDate.In(loc)
	} else {
		end = time.Now().In(loc)
	}

	totalRevenue, totalTransaction, err := s.repo.GetSummaryByPeriod(start, end)
	if err != nil {
		return nil, err
	}

	productName, productQuantitySold, err := s.repo.GetBestSellingProductByPeriod(start, end)
	if err != nil {
		return nil, err
	}

	var bestProduct *models.BestSellingProduct
	if productName != "" || productQuantitySold > 0 {
		bestProduct = &models.BestSellingProduct{
			Name:         productName,
			QuantitySold: productQuantitySold,
		}
	}

	return &models.TodayReport{
		TotalRevenue:       totalRevenue,
		TotalTransaction:   totalTransaction,
		BestSellingProduct: bestProduct,
	}, nil
}
