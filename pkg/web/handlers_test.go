package web

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/ruvcoindev/idealcore/pkg/config"
)

func TestHealthHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	cfg := config.Load()
	h := &Handlers{cfg: cfg}
	router.GET("/health", h.healthHandler)

	req, _ := http.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestGetThemesHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	cfg := config.Load()
	h := &Handlers{cfg: cfg}
	router.GET("/api/intention/themes", h.getThemesHandler)

	req, _ := http.NewRequest("GET", "/api/intention/themes", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}
