package otp

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// OTP is a struct
type OTP struct {
	Token     string
	CreatedAt time.Time
}

// OTPSvc is a struct
type OTPSvc struct {
	DB  *sync.Map
	TTL time.Duration
}

// Create is a func
func (otpSvc *OTPSvc) Create(phone string) (string, OTP) {
	rand.Seed(time.Now().UTC().UnixNano())
	token := fmt.Sprintf("%06d", rand.Intn(999999))

	otp := OTP{token, time.Now()}
	otpSvc.DB.Store(phone, otp)
	return phone, otp
}

// Verify is a func
func (otpSvc *OTPSvc) Verify(phone string, token string) bool {
	otp, ok := otpSvc.DB.Load(phone)

	if ok && otp.(OTP).Token == token {
		otpSvc.DB.Delete(phone)

		if time.Now().Sub(otp.(OTP).CreatedAt) < otpSvc.TTL {
			return true
		}

		return false
	}

	return false
}
