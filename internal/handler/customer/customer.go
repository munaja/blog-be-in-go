package customer

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/munaja/blog-practice-be-using-go/internal/handler/allbasiccommon/account"
	"github.com/munaja/blog-practice-be-using-go/internal/handler/allbasiccommon/accountom"
	"github.com/munaja/blog-practice-be-using-go/internal/handler/allbasiccommon/auth"
	er "github.com/munaja/blog-practice-be-using-go/internal/handler/allbasiccommon/errors"
	"github.com/munaja/blog-practice-be-using-go/internal/handler/allbasiccommon/home"
	"github.com/munaja/blog-practice-be-using-go/internal/handler/customer/profile"
	profileom "github.com/munaja/blog-practice-be-using-go/internal/handler/customer/profile-om"
	mu "github.com/munaja/blog-practice-be-using-go/internal/model/user"
)

func SetRoutes() http.Handler {
	auth.Position = mu.UPCustomer

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)

	r.NotFound(er.NotFoundResponse)
	r.MethodNotAllowed(er.MethodNotAllowedResponse)

	r.Get("/", home.Index)

	r.Route("/auth", func(r chi.Router) {
		r.Post("/login", auth.Login)
		r.Get("/logout", auth.Logout)
	})

	r.Route("/account", func(r chi.Router) {
		r.Post("/register", account.Register)
		r.Get("/confirm-by-email", account.ConfirmByEmail)
		r.Post("/resend-confirmation-email", account.ResendConfirmationEmail)
		r.Post("/request-reset-password", account.RequestResetPassword) // request
		r.Get("/check-reset-password", account.CheckResetPassword)      // check
		r.Patch("/reset-password", account.ResetPassword)
	})

	r.Route("/accountom", func(r chi.Router) {
		r.Use(auth.GuardMW)
		r.Patch("/change-password", accountom.ChangePassword)
		r.Get("/", accountom.GetDetail)
	})

	r.Route("/profile", func(r chi.Router) {
		r.Get("/{user_name}", profile.GetDetail)
	})

	r.Route("/profile-om", func(r chi.Router) {
		r.Use(auth.GuardMW)
		r.Get("/", profileom.GetDetail)
		r.Patch("/", profileom.Update)
	})

	return r
}
