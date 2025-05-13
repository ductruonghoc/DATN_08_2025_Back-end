package main

import (
	"github.com/ductruonghoc/DATN_08_2025_Back-end/config"
	"github.com/ductruonghoc/DATN_08_2025_Back-end/routes"
	"github.com/ductruonghoc/DATN_08_2025_Back-end/models"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	//Config env
	config.LoadEnv()

	// Try to get DSN from environment variable
	dsn := config.GetEnv("POSTGRES_DSN", "")
	if dsn == "" {
		log.Println("POSTGRES_DSN environment variable not set. Using default DSN (ensure it's configured for your local setup).")
		dsn = "postgres://user:password@localhost:5432/dbname?sslmode=disable"
	}

	DB, err := models.InitDB(dsn)
	if err != nil {
		log.Fatalf("Could not initialize database connection: %v", err)
	}
	defer DB.Close() // Ensure the database connection is closed when main exits.

	r := gin.Default()

	// Register all routes
	routes.RegisterRoutes(r)

	// Start server
	port := config.GetEnv("PORT", "8080")
	r.Run(":" + port)
}
