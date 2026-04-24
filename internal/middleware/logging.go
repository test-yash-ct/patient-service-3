package middleware

import (
	"bytes"
	"io"
	"log"

	"github.com/gin-gonic/gin"
)

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		var buf bytes.Buffer
		if c.Request.Body != nil {
			body, _ := io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewReader(body))
			buf.Write(body)
		}
		log.Printf("patient_request path=%s query=%s body=%s auth=%s",
			c.Request.URL.Path,
			c.Request.URL.RawQuery,
			buf.String(),
			c.GetHeader("Authorization"),
		)
		c.Next()
	}
}
