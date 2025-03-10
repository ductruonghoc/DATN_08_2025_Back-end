package main

import (
	"github.com/ductruonghoc/DATN_08_2025_Back-end/config"
	"github.com/ductruonghoc/DATN_08_2025_Back-end/routes"
	"github.com/gin-gonic/gin"
);

func main(){
	//Config env
	config.LoadEnv();

	r := gin.Default();

	// Register all routes
	routes.RegisterRoutes(r);

	// Start server
	port := config.GetEnv("PORT", "8080");
	r.Run(":" + port);
};