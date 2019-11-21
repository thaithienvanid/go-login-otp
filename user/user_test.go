package user

import (
	"go-login-otp/otp"
	"go-login-otp/sms"
	helper "go-login-otp/util/http"

	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type UserSvcTestSuite struct {
	suite.Suite
	userSvc UserSvc
}

func (suite *UserSvcTestSuite) SetupTest() {
	var otpDB sync.Map
	otpSvc := otp.OTPSvc{DB: &otpDB, TTL: time.Duration(60) * time.Second}

	smsSvc := sms.SMSSvc{}

	suite.userSvc = UserSvc{OTPSvc: &otpSvc, SMSSvc: &smsSvc}
}

func TestUserSvcTestSuite(t *testing.T) {
	suite.Run(t, new(UserSvcTestSuite))
}

func (suite *UserSvcTestSuite) TestIssueOTPCode() {

	// Valid request
	{
		expected := helper.Response{Message: http.StatusText(200)}

		input, _ := json.Marshal(IssueOTPCodeRequest{Phone: "0123456"})

		req := httptest.NewRequest("POST", "/", bytes.NewBuffer(input))
		w := httptest.NewRecorder()

		suite.userSvc.IssueOTPCode(w, req)

		res := w.Result()

		assert.Equal(suite.T(), 200, res.StatusCode)

		var actual helper.Response
		_ = json.NewDecoder(res.Body).Decode(&actual)

		assert.Equal(suite.T(), expected, actual)
	}

	// Invalid request
	{
		expected := helper.Response{Message: http.StatusText(400)}

		req := httptest.NewRequest("POST", "/", nil)
		w := httptest.NewRecorder()

		suite.userSvc.IssueOTPCode(w, req)

		res := w.Result()

		assert.Equal(suite.T(), 400, res.StatusCode)

		var actual helper.Response
		_ = json.NewDecoder(res.Body).Decode(&actual)

		assert.Equal(suite.T(), expected, actual)
	}

	// Invalid request
	{
		expected := helper.Response{Message: http.StatusText(400)}

		input, _ := json.Marshal(IssueOTPCodeRequest{Phone: ""})

		req := httptest.NewRequest("POST", "/", bytes.NewBuffer(input))
		w := httptest.NewRecorder()

		suite.userSvc.IssueOTPCode(w, req)

		res := w.Result()

		assert.Equal(suite.T(), 400, res.StatusCode)

		var actual helper.Response
		_ = json.NewDecoder(res.Body).Decode(&actual)

		assert.Equal(suite.T(), expected, actual)
	}

}

func (suite *UserSvcTestSuite) TestVerifyOTPCode() {

	// Valid request with correct otp
	{

		phone, otp := suite.userSvc.OTPSvc.Create("0123456")

		expected := helper.Response{Message: http.StatusText(200), Payload: "eyJwaG9uZSI6IjAxMjM0NTYifQ=="}

		input, _ := json.Marshal(VerifyOTPCodeRequest{Phone: phone, Token: otp.Token})

		req := httptest.NewRequest("POST", "/", bytes.NewBuffer(input))
		w := httptest.NewRecorder()

		suite.userSvc.VerifyOTPCode(w, req)

		res := w.Result()

		assert.Equal(suite.T(), 200, res.StatusCode)

		var actual helper.Response
		_ = json.NewDecoder(res.Body).Decode(&actual)

		assert.Equal(suite.T(), expected, actual)
	}

	// Valid request with incorrect otp
	{

		phone, _ := suite.userSvc.OTPSvc.Create("0123456")

		expected := helper.Response{Message: http.StatusText(401)}

		input, _ := json.Marshal(VerifyOTPCodeRequest{Phone: phone, Token: "xxxxxx"})

		req := httptest.NewRequest("POST", "/", bytes.NewBuffer(input))
		w := httptest.NewRecorder()

		suite.userSvc.VerifyOTPCode(w, req)

		res := w.Result()

		assert.Equal(suite.T(), 401, res.StatusCode)

		var actual helper.Response
		_ = json.NewDecoder(res.Body).Decode(&actual)

		assert.Equal(suite.T(), expected, actual)
	}

	// Invalid request
	{
		expected := helper.Response{Message: http.StatusText(400)}

		// input, _ := json.Marshal(VerifyOTPCodeRequest{Phone: "", Token: ""})

		req := httptest.NewRequest("POST", "/", nil)
		w := httptest.NewRecorder()

		suite.userSvc.VerifyOTPCode(w, req)

		res := w.Result()

		assert.Equal(suite.T(), 400, res.StatusCode)

		var actual helper.Response
		_ = json.NewDecoder(res.Body).Decode(&actual)

		assert.Equal(suite.T(), expected, actual)
	}

}
