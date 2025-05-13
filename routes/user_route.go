// For sub route groups
package routes

import (
	"database/sql"
	"github.com/ductruonghoc/DATN_08_2025_Back-end/controllers"
	"github.com/ductruonghoc/DATN_08_2025_Back-end/middlewares"
	"github.com/gin-gonic/gin"
)

func UserRoutes(r *gin.Engine, db *sql.DB) {
	userGroup := r.Group("/user")
	{
		userGroup.POST(
			"/nonverified_registration",
			middlewares.CheckVerifiedEmailExisted(db),
			middlewares.UserExistedIgnore(),
			middlewares.StoreTemporatoryUser(db),
			middlewares.SendOTP(),
			controllers.NonVerifiedRegistration,
		);
		userGroup.POST(
			"/verified_registration",
			middlewares.VeirifyOTP(),
			controllers.VerifiedRegistration(),
		);
		userGroup.POST(
			"/can_google_register",
			middlewares.CheckGoogleUserExisted(),
			middlewares.UserExistedIgnore(),
			controllers.CanGoogleRegister,
		);
		userGroup.POST(
			"/google_registration",
			controllers.GoogleRegistration,
		)
		userGroup.POST(
			"/login",
			middlewares.UserAuthenticate(),
			controllers.UserLogin,
		);
		userGroup.POST(
			"/google_login",
			middlewares.CheckGoogleUserExisted(),
			controllers.GoogleLogin,
		);
		userGroup.POST(
			"/can_reset_password",
			middlewares.CheckVerifiedEmailExisted(db),
			middlewares.UserExistedFirst(),
			middlewares.SendOTP(),
			controllers.CanResetPassword,
		);
		userGroup.POST(
			"/reset_password",
			middlewares.VeirifyOTP(),
			controllers.ResetPassword(),
		);
	}
}
