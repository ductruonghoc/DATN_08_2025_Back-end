package routes

import (
	"github.com/gin-gonic/gin"
);

func RegisterRoutes(r *gin.Engine) {
	// Add route groups here
	TemplateRoutes(r);
};
