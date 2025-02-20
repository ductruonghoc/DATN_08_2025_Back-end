package main;

import (
	"github.com/gin-gonic/gin"
	"github.com/ductruonghoc/DATN_08_2025_Back-end/routes"
);

func main(){
	r := gin.Default();

	// Register all routes
	routes.RegisterRoutes(r);

	// Start server
	r.Run(":8080");
};