package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/healthops/patient-service/internal/auth"
)

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		log.Printf("http method=%s path=%s status=%d dur_ms=%d",
			c.Request.Method,
			c.FullPath(),
			c.Writer.Status(),
			time.Since(start).Milliseconds(),
		)
	}
}

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.GetHeader("X-Request-ID")
		if id == "" || len(id) > 128 {
			id = randomRequestID()
		}
		c.Set("request_id", id)
		c.Header("X-Request-ID", id)
		c.Next()
	}
}

func Authenticate(secret string, maxTTLSec int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, err := auth.ParseBearer(c.GetHeader("Authorization"), secret, time.Duration(maxTTLSec)*time.Second)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": "unauthorized"})
			return
		}
		c.Set("claims", claims)
		c.Next()
	}
}

func ClaimsFrom(c *gin.Context) (auth.Claims, bool) {
	v, ok := c.Get("claims")
	if !ok {
		return auth.Claims{}, false
	}
	cl, ok := v.(auth.Claims)
	return cl, ok
}
