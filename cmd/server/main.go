package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/healthops/patient-service/internal/audit"
	"github.com/healthops/patient-service/internal/config"
	"github.com/healthops/patient-service/internal/crypto"
	"github.com/healthops/patient-service/internal/dbpool"
	"github.com/healthops/patient-service/internal/handlers"
	"github.com/healthops/patient-service/internal/middleware"
	"github.com/healthops/patient-service/internal/store"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}
	enc, err := crypto.NewFieldEncryptor(cfg.PIIKey)
	if err != nil {
		log.Fatalf("crypto: %v", err)
	}
	ctx := context.Background()
	pool, err := dbpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("db: %v", err)
	}
	defer pool.Close()

	gin.SetMode(gin.ReleaseMode)
	if cfg.Debug {
		log.Printf("debug logging enabled; PHI debug surfaces are not registered")
	}

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.RequestID())
	r.Use(middleware.CORS(cfg.CORSOrigins))
	r.Use(middleware.RequestLogger())

	r.GET("/healthz", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})
	r.GET("/readyz", func(c *gin.Context) {
		if err := pool.Ping(c.Request.Context()); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "db_unavailable"})
			return
		}
		c.String(http.StatusOK, "ok")
	})

	v1 := r.Group("/v1")
	v1.Use(middleware.Authenticate(cfg.JWTSecret, cfg.MaxTokenTTLSec))
	v1.Use(middleware.RateLimit(60, time.Minute))
	papi := &handlers.PatientAPI{Store: store.New(pool, enc), Audit: audit.New()}
	papi.Register(v1)

	srv := &http.Server{
		Addr:              cfg.ListenAddr,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      60 * time.Second,
		IdleTimeout:       120 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}
	go func() {
		log.Printf("listening on %s", cfg.ListenAddr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server: %v", err)
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = srv.Shutdown(shutdownCtx)
}
