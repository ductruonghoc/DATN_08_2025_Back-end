package middlewares

import (
	"net/http"

	//"time"

	"github.com/ductruonghoc/DATN_08_2025_Back-end/internal"
	"github.com/ductruonghoc/DATN_08_2025_Back-end/models"
	"github.com/gin-gonic/gin"
	//"github.com/gin-gonic/gin/binding"
)

func CheckVerifiedEmailExisted() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Email string `json:"email"`
		}
		// Bind JSON, form, or query parameter
		if err := c.ShouldBindBodyWithJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
			c.Abort()
			return
		}
		// Blank Email request
		if req.Email == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Email is required"})
			c.Abort()
			return
		}

		account_existed := false
		//db query here
		

		c.Set("account_existed", account_existed)
		c.Next()
	}
}

func StoreTemporatoryUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		// Bind JSON, form, or query parameter
		if err := c.ShouldBindBodyWithJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
			c.Abort()
			return
		}

		//password hashing
		hashed_password := internal.BcryptHashing(req.Password)
		req.Password = hashed_password

		//db query here

		// Successfully stored unverified user, continue processing
		c.Next()
	}
}

func SendOTP() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Email string `json:"email"`
		}

		// Bind JSON, form, or query parameter
		if err := c.ShouldBindBodyWithJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
			c.Abort()
			return
		}

		if req.Email == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Email is required"})
			c.Abort()
			return
		}

		var otp models.OTP

		otpCode := internal.Digit6Random()
		//expiration := time.Now().Add(5 * time.Minute);

		otp.OTPCode = otpCode
		//otp.OTPExpiresAt = expiration;

		//db query here

		//email otp
		if err := internal.EmailOTP(req.Email, otp.OTPCode); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send OTP"})
			c.Abort()
			return
		}

		//succesful
		c.Next()
	}
}

func VeirifyOTP() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Email   string `json:"email"`
			OTPCode string `json:"otp_code"`
		}
		//try bind the request
		if err := c.ShouldBindBodyWithJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
			c.Abort()
			return
		}
		//db query here

		//Successful
		c.Set("verified_email", req.Email)
		c.Next()
	}
}

func CheckGoogleUserExisted() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			GoogleID string `json:"google_id"`
		}
		//try binding the request
		if err := c.ShouldBindBodyWithJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
			c.Abort()
			return
		}

		account_existed := false
		//db query here
		c.Set("account_existed", account_existed)
		c.Next()
	}
}

func UserExistedIgnore() gin.HandlerFunc {
	return func(c *gin.Context) {
		//Get middlewares results
		account_existed, exists := c.Get("account_existed")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve account existed status"})
			return
		}

		//account existed before will be ignored
		if account_existed == true {
			c.JSON(http.StatusConflict, gin.H{"error": "Account has already existed"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// Authenticate Middleware authenticates users based on username and password
func UserAuthenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		//password hashing
		hashed_password := internal.BcryptHashing(req.Password)
		req.Password = hashed_password

		var userID int
		//db query here
		c.Set("user_id", userID)
		c.Next()
	}
}

func UserExistedFirst() gin.HandlerFunc {
	return func(c *gin.Context) {
		//Get middlewares results
		account_existed, exists := c.Get("account_existed")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve account existed status"})
			return
		}

		//account existed before will be ignored
		if account_existed == false {
			c.JSON(http.StatusConflict, gin.H{"error": "Account has not existed yet."})
			c.Abort()
			return
		}

		c.Next()
	}
}
