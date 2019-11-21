package otp

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type OTPSvcTestSuite struct {
	suite.Suite
	otpSvc OTPSvc
}

var sm sync.Map

func (suite *OTPSvcTestSuite) SetupTest() {
	suite.otpSvc = OTPSvc{DB: &sm, TTL: time.Duration(1) * time.Second}
}

func TestOTPSvcTestSuite(t *testing.T) {
	suite.Run(t, new(OTPSvcTestSuite))
}

func (suite *OTPSvcTestSuite) TestCreate() {
	phone, value := suite.otpSvc.Create("0123456")
	assert.Equal(suite.T(), "0123456", phone)
	assert.NotEqual(suite.T(), "", value.Token)
}

func (suite *OTPSvcTestSuite) TestVerify() {

	// Valid token and valid time
	{
		_, value := suite.otpSvc.Create("0123456")
		assert.Equal(suite.T(), true, suite.otpSvc.Verify("0123456", value.Token))
	}

	// Valid token and invalid time
	{
		_, value := suite.otpSvc.Create("0123456")
		time.Sleep(suite.otpSvc.TTL)
		assert.Equal(suite.T(), false, suite.otpSvc.Verify("0123456", value.Token))
	}

	// Invalid token
	{
		_, _ = suite.otpSvc.Create("0123456")
		assert.Equal(suite.T(), false, suite.otpSvc.Verify("0123456", "xxxxxx"))
	}
}
