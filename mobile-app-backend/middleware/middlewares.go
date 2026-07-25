package middleware

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"

	"mobile-app-backend/store"

	"github.com/gin-gonic/gin"
)

type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

// Write intercepts the response data and writes it to the buffer
func (r responseBodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "unauthorized",
				"message": "Missing or invalid Authorization header",
			})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		userID, err := store.VerifyToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "unauthorized",
				"message": "Invalid or expired token: " + err.Error(),
			})
			return
		}

		c.Set("user_id", userID)
		c.Set("token", tokenString)
		c.Next()
	}
}

func LogResponse() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Initialize custom writer
		w := &responseBodyWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = w

		// 2. Process the request handler
		c.Next()

		reqMethod := c.Request.Method

		// 3. Handle the captured response
		statusCode := c.Writer.Status()
		responseBody := w.body.String()
		fmt.Printf("\n\n[Respone] Status: %d | Body: %s\n", statusCode, responseBody)
		fmt.Printf("[Request] Method: %s | Path: %s\n\n", reqMethod, c.Request.RequestURI)
	}
}
