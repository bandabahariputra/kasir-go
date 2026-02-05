package main

import (
	"encoding/json"
	"fmt"
	"kasir-go/database"
	"kasir-go/handlers"
	"kasir-go/repositories"
	"kasir-go/services"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Port   string `mapstructure:"PORT"`
	DBConn string `mapstructure:"DB_CONN"`
}

func main() {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		_ = viper.ReadInConfig()
	}

	config := Config{
		Port:   viper.GetString("PORT"),
		DBConn: viper.GetString("DB_CONN"),
	}

	// setup database
	db, err := database.InitDB(config.DBConn)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	categoryRepo := repositories.NewCategoryRepository(db)
	productRepo := repositories.NewProductRepository(db)
	transactionRepo := repositories.NewTransactionRepository(db)

	categoryService := services.NewCategoryService(categoryRepo, productRepo)
	productService := services.NewProductService(productRepo, categoryRepo)
	transactionService := services.NewTransactionService(transactionRepo)
	reportService := services.NewReportService(transactionRepo)

	categoryHandler := handlers.NewCategoryHandler(categoryService)
	productHandler := handlers.NewProductHandler(productService)
	transactionHandler := handlers.NewTransactionHandler(transactionService)
	reportHandler := handlers.NewReportHandler(reportService)

	http.HandleFunc("/api/categories/", categoryHandler.HandleCategoryByID)
	http.HandleFunc("/api/categories", categoryHandler.HandleCategories)

	http.HandleFunc("/api/products/", productHandler.HandleProductByID)
	http.HandleFunc("/api/products", productHandler.HandleProducts)

	http.HandleFunc("/api/checkout", transactionHandler.Checkout)

	http.HandleFunc("/api/report/today", reportHandler.GetTodayReport)

	// GET http://localhost:8080/health
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status": "up",
		})
	})

	fmt.Println("Server starting on port " + config.Port + "...")

	err = http.ListenAndServe(":"+config.Port, nil)

	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
