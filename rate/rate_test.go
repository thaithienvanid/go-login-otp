package rate

import (
	helper "go-login-otp/util/http"
	"time"

	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	timerate "golang.org/x/time/rate"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type RateSvcTestSuite struct {
	suite.Suite
	rateSvc RateSvc
}

var sm sync.Map

func (suite *RateSvcTestSuite) SetupTest() {
	suite.rateSvc = RateSvc{DB: &sm, Limit: timerate.Every(1 * time.Second), Burst: 1}
}

func TestRateSvcTestSuite(t *testing.T) {
	suite.Run(t, new(RateSvcTestSuite))
}

func (suite *RateSvcTestSuite) TestRate() {

	// Valid address and valid rate
	{
		expected := helper.Response{Message: http.StatusText(200)}

		fn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			helper.JSON(w, http.StatusOK, expected)
		})

		handler := suite.rateSvc.Rate(fn)

		req := httptest.NewRequest("GET", "/", nil)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		res := w.Result()

		assert.Equal(suite.T(), 200, res.StatusCode)

		var actual helper.Response
		_ = json.NewDecoder(res.Body).Decode(&actual)

		assert.Equal(suite.T(), expected, actual)
	}

	// Valid address and invalid rate
	{
		expected := helper.Response{Message: http.StatusText(429)}

		fn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			helper.JSON(w, http.StatusOK, helper.Response{Message: http.StatusText(200)})
		})

		handler := suite.rateSvc.Rate(fn)

		req := httptest.NewRequest("GET", "/", nil)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		res := w.Result()

		assert.Equal(suite.T(), 429, res.StatusCode)

		var actual helper.Response
		_ = json.NewDecoder(res.Body).Decode(&actual)

		assert.Equal(suite.T(), expected, actual)
	}

	// Invalid address
	{
		expected := helper.Response{Message: http.StatusText(500)}

		fn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			helper.JSON(w, http.StatusOK, helper.Response{Message: http.StatusText(200)})
		})

		handler := suite.rateSvc.Rate(fn)

		req := httptest.NewRequest("GET", "/", nil)
		req.RemoteAddr = "invalid"
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		res := w.Result()

		assert.Equal(suite.T(), 500, res.StatusCode)

		var actual helper.Response
		_ = json.NewDecoder(res.Body).Decode(&actual)

		assert.Equal(suite.T(), expected, actual)
	}

}
