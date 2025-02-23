package controllers

import (
	"net/http"

	"github.com/ductruonghoc/DATN_08_2025_Back-end/internal"
	"github.com/ductruonghoc/DATN_08_2025_Back-end/models"
	"github.com/gin-gonic/gin"
)

func NonVerifiedRegistration(c *gin.Context) {
	//all middleware succesfully processes
	c.JSON(http.StatusOK, gin.H{"message": "Nonverified registration process completed successfully"})
}

func VerifiedRegistration() gin.HandlerFunc {
	return func(c *gin.Context) {
		//Get middlewares results
		email, exists := c.Get("verified_email")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve verified email"})
			return
		}

		//db query here

		//Successful
		c.JSON(http.StatusOK, gin.H{"message": "User verified successfully", "email": email})
	}
}

func CanGoogleRegister(c *gin.Context) {
	//all middleware succesfully processes
	c.JSON(http.StatusOK, gin.H{"message": "Google Registration can process"})
}

func GoogleRegistration(c *gin.Context) {
	var req struct {
		GoogleID    string `json:"google_id"`
		GoogleEmail string `json:"google_email"`
	}

	// Bind JSON, form, or query parameter
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		c.Abort()
		return
	}

	user := models.User{
		Email: "",
	}

	id, err := models.InsertUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
		c.Abort()
		return
	}

	google_user := models.GoogleUser{
		GoogleID:    req.GoogleID,
		GoogleEmail: req.GoogleEmail,
	}

	if _, err := models.InsertUserGoogleInfomation(google_user, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
		c.Abort()
		return
	}

	//succesfully process
	c.JSON(http.StatusOK, gin.H{"message": "Google Registration process successfully"})
}

func UserLogin(c *gin.Context) {
	//Get middlewares results
	payload, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve account existed status"})
		return
	}

	userID, ok := payload.(int)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	token, err := internal.JWTGenerator(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func GoogleLogin(c *gin.Context) {
	//Get middlewares results
	account_existed, exists := c.Get("account_existed")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve account existed status"})
		return
	}
	//not existed yet
	if account_existed == false {
		c.JSON(http.StatusNotFound, gin.H{"error": "Could not retrieve account"})
		return
	}
	//db query here
	var userID int
	token, err := internal.JWTGenerator(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func CanResetPassword(c *gin.Context) {
	//all middleware succesfully processes
	c.JSON(http.StatusOK, gin.H{"message": "Reset Password can process. OTP has been sent"})
}

func ResetPassword() gin.HandlerFunc{
	return func(c *gin.Context) {
		//Get middlewares results
		_, exists := c.Get("verified_email");
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve verified email"});
			return;
		}

		//db query here

		//Successful
		c.JSON(http.StatusOK, gin.H{"message": "Password resets successfully."});
	}
}
