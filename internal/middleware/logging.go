package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/healthops/patient-service/internal/obs"
)

// RequestLogger emits structured JSON access logs via the obs package.
func RequestLogger() gin.HandlerFunc {
	return obs.AccessLogger()
}
