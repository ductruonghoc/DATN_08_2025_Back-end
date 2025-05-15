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
			"/unverified_register",
			middlewares.CheckVerifiedEmailExisted(db),
			middlewares.UserExistedIgnore(),
			middlewares.SendOTP(),
			middlewares.StoreTemporatoryUser(db),
			controllers.NonVerifiedRegistration,
		);
		userGroup.POST(
			"/verify_registration",
			middlewares.VeirifyOTP_Register(db),
			controllers.VerifiedRegistration(db),
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
			middlewares.UserAuthenticate(db),
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
			middlewares.VeirifyOTP_Register(db),
			controllers.ResetPassword(),
		);
	}
}
