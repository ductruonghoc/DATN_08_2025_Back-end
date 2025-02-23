package models;

import (
	"time"
);

type User struct {
	Email string `json:"email"`
};

type OTP struct {
	OTPCode string `json:"otp_code"`
	OTPExpiresAt time.Time `json:"expires_at"`
};

type GoogleUser struct {
	GoogleID string `json:"google_id"`
	GoogleEmail string `json:"google_email"`
};