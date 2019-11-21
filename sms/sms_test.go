package sms

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type SMSSvcTestSuite struct {
	suite.Suite
	smsSvc SMSSvc
}

func (suite *SMSSvcTestSuite) SetupTest() {
	suite.smsSvc = SMSSvc{}
}

func TestSMSSvcTestSuite(t *testing.T) {
	suite.Run(t, new(SMSSvcTestSuite))
}

func (suite *SMSSvcTestSuite) TestSend() {
	assert.Equal(suite.T(), nil, suite.smsSvc.Send("0123456789", "000000"))
}
