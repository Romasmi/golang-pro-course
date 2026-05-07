package internalhttp

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockLogger struct {
	mock.Mock
}

func (m *mockLogger) Info(msg string)  { m.Called(msg) }
func (m *mockLogger) Error(msg string) { m.Called(msg) }

func TestLoggingMiddleware(t *testing.T) {
	logger := &mockLogger{}
	logger.On("Info", mock.Anything).Return()

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusAccepted)
		_, _ = w.Write([]byte("OK"))
	})

	handler := loggingMiddleware(logger, nextHandler)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusAccepted, rr.Code)
	require.Equal(t, "OK", rr.Body.String())
	logger.AssertCalled(t, "Info", mock.Anything)
}
