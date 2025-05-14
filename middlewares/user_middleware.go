package middlewares

import (
	"database/sql"
	"net/http"

	"time"

	"github.com/ductruonghoc/DATN_08_2025_Back-end/internal"
	"github.com/ductruonghoc/DATN_08_2025_Back-end/models"

	"golang.org/x/crypto/bcrypt"
	"github.com/gin-gonic/gin"
	//"github.com/gin-gonic/gin/binding"
)

func CheckVerifiedEmailExisted(db *sql.DB) gin.HandlerFunc {
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
		email := req.Email
		account_existed := false
		//db query here
		query := `
			select 1 as account_existed
			from "user" 
			where email = $1
			limit 1
		`
		rows, err := db.Query(query, email) // Using a placeholder for the argument
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Can't query"})
			c.Abort()
			return
		}
		defer rows.Close() // Important to close rows to free resources

		if rows.Next() {
			// A row exists, which means the account exists
			account_existed = true
		}

		c.Set("account_existed", account_existed)
		c.Next()
	}
}

func StoreTemporatoryUser(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		//Get middlewares results
		otpVal, exists := c.Get("otp")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve otp existed status"})
			c.Abort()
			return
		}

		otp, ok := otpVal.(models.OTP) // type assertion
		if !ok {
			// handle type mismatch
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error with otp"})
			c.Abort()
			return
		}

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
		hashed_otp := internal.BcryptHashing(otp.OTPCode)
		req.Password = hashed_password

		email := req.Email
		//db query here
		query := `
			INSERT INTO temp_user (email, password, otp, otp_generated_time)
			VALUES ($1, $2, $3, $4)
			ON CONFLICT (email) DO UPDATE
			SET password = EXCLUDED.password,
				otp = EXCLUDED.otp,
				otp_generated_time = EXCLUDED.otp_generated_time;
		`
		rows, err := db.Query(query, email, hashed_password, hashed_otp, otp.OTPWasGeneratedAt) // Using a placeholder for the argument
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Can't query"})
			c.Abort()
			return
		}
		defer rows.Close() // Important to close rows to free resources
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
		otp.OTPWasGeneratedAt = time.Now()

		//email otp
		if err := internal.EmailOTP(req.Email, otp.OTPCode); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send OTP"})
			c.Abort()
			return
		}
		c.Set("otp", otp)
		//succesful
		c.Next()
	}
}

func VeirifyOTP_Register(db *sql.DB) gin.HandlerFunc {
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
		email := req.Email
		query := `
			select 
				otp as otp_code,
				otp_generated_time as generated_at
			from temp_user
			where email = $1;
		`
		rows, err := db.Query(query, email) // Using a placeholder for the argument
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Can't query"})
			c.Abort()
			return
		}
		defer rows.Close() // Important to close rows to free resources

		var otp models.OTP
		if rows.Next() { // Advances to the first (and expected only) row.
			// Could return false here if there's an immediate error fetching the first row.
			if err := rows.Scan(&otp.OTPCode, &otp.OTPWasGeneratedAt); err != nil {
				// This error is specific to scanning THIS row's data.
				// e.g., otp_code was NULL and OTPCode is a non-pointer string, or types are incompatible.
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process data"})
				c.Abort()
				return
			}
			// If Scan was successful, otp is populated.
		} else {
			// rows.Next() returned false. This could be because:
			// 1. No rows were found (sql.ErrNoRows if using QueryRow().Scan(), but with raw Next() it's just 'false').
			// 2. An error occurred trying to fetch the first row.
			// This 'else' block in your code assumes it's "No rows found".
			c.JSON(http.StatusNotFound, gin.H{"error": "OTP not found"})
			c.Abort()
			return
		}

		// Check for errors after the potential Next() call
		// This is where you catch the error if rows.Next() returned 'false' due to an error,
		// rather than just no rows.
		if err = rows.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error processing results"})
			c.Abort()
			return
		}

		// Get the current time
		currentTime := time.Now()

		// Add 2 hours to the input time
		expirationTime := otp.OTPWasGeneratedAt.Add(2 * time.Hour)
		otp_is_expired := expirationTime.Before(currentTime)

		if otp_is_expired {
			c.JSON(http.StatusBadRequest, gin.H{"error": "OTP expired"})
			c.Abort()
			return
		}

		// Compare the hashed password with the plain text one
		err = bcrypt.CompareHashAndPassword([]byte(otp.OTPCode), []byte(req.OTPCode))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "OTP is not match"})
			c.Abort()
			return
		}

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
