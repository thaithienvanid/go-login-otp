package rate

import (
	helper "go-login-otp/util/http"

	"net"
	"net/http"
	"sync"

	timerate "golang.org/x/time/rate"
)

// RateSvc is a struct
type RateSvc struct {
	DB    *sync.Map
	Limit timerate.Limit
	Burst int
}

// Rate is a func
func (rateSvc *RateSvc) Rate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		host, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			helper.JSON(w, http.StatusInternalServerError, helper.Response{Message: http.StatusText(500)})
			return
		}

		limiter, exists := rateSvc.DB.Load(host)

		if !exists {
			limiter = timerate.NewLimiter(rateSvc.Limit, rateSvc.Burst)
			rateSvc.DB.Store(host, limiter)
		}

		if limiter.(*timerate.Limiter).Allow() == false {
			helper.JSON(w, http.StatusTooManyRequests, helper.Response{Message: http.StatusText(429)})
			return
		}

		next.ServeHTTP(w, r)
	})
}
