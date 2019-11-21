package user

import (
	"go-login-otp/otp"
	"go-login-otp/sms"
	helper "go-login-otp/util/http"

	"encoding/base64"
	"encoding/json"
	"net/http"
)

// UserSvc is a struct
type UserSvc struct {
	OTPSvc *otp.OTPSvc
	SMSSvc *sms.SMSSvc
}

// IssueOTPCodeRequest is a struct
type IssueOTPCodeRequest struct {
	Phone string `json:"phone,omitempty"`
}

// IssueOTPCode is a func
func (userSvc *UserSvc) IssueOTPCode(w http.ResponseWriter, r *http.Request) {
	var body IssueOTPCodeRequest

	err := json.NewDecoder(r.Body).Decode(&body)

	if err != nil || body.Phone == "" {
		helper.JSON(w, http.StatusBadRequest, helper.Response{Message: http.StatusText(400)})
		return
	}

	_, value := userSvc.OTPSvc.Create(body.Phone)
	userSvc.SMSSvc.Send(body.Phone, value.Token)

	helper.JSON(w, http.StatusOK, helper.Response{Message: http.StatusText(200)})
}

// VerifyOTPCodeRequest is a struct
type VerifyOTPCodeRequest struct {
	Phone string `json:"phone,omitempty"`
	Token string `json:"token,omitempty"`
}

// VerifyOTPCode is a func
func (userSvc *UserSvc) VerifyOTPCode(w http.ResponseWriter, r *http.Request) {
	var body VerifyOTPCodeRequest

	err := json.NewDecoder(r.Body).Decode(&body)

	if err != nil || body.Phone == "" || body.Token == "" {
		helper.JSON(w, http.StatusBadRequest, helper.Response{Message: http.StatusText(400)})
		return
	}

	res := userSvc.OTPSvc.Verify(body.Phone, body.Token)

	if res == false {
		helper.JSON(w, http.StatusUnauthorized, helper.Response{Message: http.StatusText(401)})
		return
	}

	// Just a simple token which never expire
	type Secret struct {
		Phone string `json:"phone,omitempty"`
	}
	input, _ := json.Marshal(Secret{Phone: body.Phone})
	token := base64.URLEncoding.EncodeToString(input)

	helper.JSON(w, http.StatusOK, helper.Response{Message: http.StatusText(200), Payload: token})
}
