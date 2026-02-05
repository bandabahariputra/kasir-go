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

	var productName string
	var productQuantitySold int
	productName, productQuantitySold, err = s.repo.GetBestSellingProductByPeriod(start, end)
	if err != nil {
		return nil, err
	}

	return &models.TodayReport{
		TotalRevenue:     totalRevenue,
		TotalTransaction: totalTransaction,
		BestSellingProduct: &models.BestSellingProduct{
			Name:         productName,
			QuantitySold: productQuantitySold,
		},
	}, err
}
