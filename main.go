package main

import (
	"go-login-otp/otp"
	"go-login-otp/rate"
	"go-login-otp/sms"
	"go-login-otp/user"
	helper "go-login-otp/util/http"

	"log"
	"net/http"
	"sync"
	"time"

	timerate "golang.org/x/time/rate"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func main() {

	router := chi.NewRouter()

	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	var otpDB sync.Map
	otpSvc := otp.OTPSvc{DB: &otpDB, TTL: time.Duration(60) * time.Second}

	smsSvc := sms.SMSSvc{}

	userSvc := user.UserSvc{OTPSvc: &otpSvc, SMSSvc: &smsSvc}

	var rateDB sync.Map
	rateSvc := rate.RateSvc{DB: &rateDB, Limit: timerate.Every(1 * time.Second), Burst: 1}

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		helper.Text(w, http.StatusOK, "Welcome")
	})

	router.Route("/user/login/phone", func(r chi.Router) {
		r.Use(rateSvc.Rate)
		r.Post("/", userSvc.IssueOTPCode)
	})

	router.Post("/user/login/phone/callback", userSvc.VerifyOTPCode)

	log.Println("Application is running on port 3000")
	http.ListenAndServe(":3000", router)
}
