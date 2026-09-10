package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/healthops/patient-service/internal/handlers"
	"github.com/healthops/patient-service/internal/obs"
	"github.com/healthops/patient-service/internal/store"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(obs.RequestID())
	r.Use(obs.AccessLogger())
	r.GET("/healthz", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})
	r.GET("/meta", obs.MetaHandler(obs.Metadata{
		Service:   "patient-service",
		Version:   "test",
		BuildTime: "2026-01-01T00:00:00Z",
		GitSHA:    "abc123",
	}))
	return r
}

func TestHealthRoute(t *testing.T) {
	r := setupRouter()
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status %d", w.Code)
	}
}

func TestMetaEndpoint(t *testing.T) {
	r := setupRouter()
	req := httptest.NewRequest(http.MethodGet, "/meta", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status %d", w.Code)
	}
	var body map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body["service"] != "patient-service" {
		t.Fatalf("service %q", body["service"])
	}
	if body["version"] != "test" {
		t.Fatalf("version %q", body["version"])
	}
}

func TestRequestIDEcho(t *testing.T) {
	r := setupRouter()
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	req.Header.Set(obs.HeaderRequestID, "req-test-123")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if got := w.Header().Get(obs.HeaderRequestID); got != "req-test-123" {
		t.Fatalf("expected echoed request id, got %q", got)
	}
}

func TestRequestIDGeneratedWhenMissing(t *testing.T) {
	r := setupRouter()
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if got := w.Header().Get(obs.HeaderRequestID); got == "" {
		t.Fatal("expected generated request id")
	}
}

func TestPatientHandlerRequiresAuth(t *testing.T) {
	r := setupRouter()
	api := &handlers.PatientAPI{Store: (*store.PatientStore)(nil), Secret: "x"}
	g := r.Group("/v1")
	api.Register(g)
	req := httptest.NewRequest(http.MethodGet, "/v1/patients/p-1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 got %d", w.Code)
	}
}

func TestMetadataFromEnvDefaults(t *testing.T) {
	os.Unsetenv("SERVICE_VERSION")
	os.Unsetenv("GIT_SHA")
	os.Unsetenv("BUILD_TIME")
	meta := obs.MetadataFromEnv("patient-service")
	if meta.Version != "dev" || meta.GitSHA != "unknown" || meta.BuildTime != "unknown" {
		t.Fatalf("unexpected defaults: %+v", meta)
	}
}
