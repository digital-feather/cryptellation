package health

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestHealthSuite(t *testing.T) {
	suite.Run(t, new(HealthSuite))
}

type HealthSuite struct {
	suite.Suite
}

func (suite *HealthSuite) TestLiveness() {
	h := New()
	req := httptest.NewRequest(http.MethodGet, "/liveness", nil)
	w := httptest.NewRecorder()
	livenessFunc := h.liveness()

	livenessFunc(w, req)
	suite.Require().Equal(http.StatusOK, w.Result().StatusCode)
}

func (suite *HealthSuite) TestReadiness() {
	h := New()
	req := httptest.NewRequest(http.MethodGet, "/readiness", nil)
	readinessFunc := h.readiness()

	h.Ready(false)
	w := httptest.NewRecorder()
	readinessFunc(w, req)
	suite.Require().Equal(http.StatusServiceUnavailable, w.Result().StatusCode)

	h.Ready(true)
	w = httptest.NewRecorder()
	readinessFunc(w, req)
	suite.Require().Equal(http.StatusOK, w.Result().StatusCode)

	h.Ready(false)
	w = httptest.NewRecorder()
	readinessFunc(w, req)
	suite.Require().Equal(http.StatusServiceUnavailable, w.Result().StatusCode)
}
