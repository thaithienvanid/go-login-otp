package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type HTTPTestSuite struct {
	suite.Suite
}

func (suite *HTTPTestSuite) SetupTest() {

}

func TestHTTPTestSuite(t *testing.T) {
	suite.Run(t, new(HTTPTestSuite))
}

func (suite *HTTPTestSuite) TestJSON() {
	expected := Response{Message: http.StatusText(200)}

	handler := func(w http.ResponseWriter, r *http.Request) {
		JSON(w, http.StatusOK, expected)
	}

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	res := w.Result()

	assert.Equal(suite.T(), 200, res.StatusCode)

	var actual Response
	_ = json.NewDecoder(res.Body).Decode(&actual)

	assert.Equal(suite.T(), expected, actual)
}

func (suite *HTTPTestSuite) TestText() {
	expected := "OK"

	handler := func(w http.ResponseWriter, r *http.Request) {
		Text(w, http.StatusOK, expected)
	}

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	res := w.Result()

	assert.Equal(suite.T(), 200, res.StatusCode)

	var actual string
	_ = json.NewDecoder(res.Body).Decode(&actual)

	assert.Equal(suite.T(), expected, actual)
}
