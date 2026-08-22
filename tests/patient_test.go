package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/healthops/patient-service/internal/handlers"
	"github.com/healthops/patient-service/internal/middleware"
	"github.com/healthops/patient-service/internal/store"
)

func TestHealthRoute(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/healthz", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status %d", w.Code)
	}
}

func TestPatientHandlerRequiresAuth(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	api := &handlers.PatientAPI{Store: (*store.PatientStore)(nil)}
	g := r.Group("/v1")
	g.Use(middleware.Authenticate("x", 900))
	api.Register(g)
	req := httptest.NewRequest(http.MethodGet, "/v1/patients/p-1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 got %d", w.Code)
	}
}
