package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"mobile-app-backend/store"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func newTestRouter() *gin.Engine {
	return gin.New()
}

func TestCORSMiddleware_SetsHeaders(t *testing.T) {
	router := newTestRouter()
	router.Use(CORSMiddleware())
	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
	if got := w.Header().Get("Access-Control-Allow-Origin"); got != "*" {
		t.Errorf("expected Access-Control-Allow-Origin '*', got %q", got)
	}
	if got := w.Header().Get("Access-Control-Allow-Credentials"); got != "true" {
		t.Errorf("expected Access-Control-Allow-Credentials 'true', got %q", got)
	}
	if got := w.Header().Get("Access-Control-Allow-Methods"); got == "" {
		t.Errorf("expected Access-Control-Allow-Methods header to be set")
	}
}

func TestCORSMiddleware_OptionsRequestShortCircuits(t *testing.T) {
	router := newTestRouter()
	router.Use(CORSMiddleware())
	handlerCalled := false
	router.OPTIONS("/ping", func(c *gin.Context) {
		handlerCalled = true
		c.String(http.StatusOK, "pong")
	})

	req := httptest.NewRequest(http.MethodOptions, "/ping", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected status 204, got %d", w.Code)
	}
	if handlerCalled {
		t.Error("expected downstream handler not to be called for OPTIONS request")
	}
}

func TestAuthMiddleware_MissingHeader(t *testing.T) {
	router := newTestRouter()
	router.Use(AuthMiddleware())
	router.GET("/protected", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", w.Code)
	}
}

func TestAuthMiddleware_HeaderWithoutBearerPrefix(t *testing.T) {
	router := newTestRouter()
	router.Use(AuthMiddleware())
	router.GET("/protected", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "sometoken")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", w.Code)
	}
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	router := newTestRouter()
	router.Use(AuthMiddleware())
	router.GET("/protected", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer not-a-real-token")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", w.Code)
	}
}

func TestAuthMiddleware_ValidToken(t *testing.T) {
	userID := uuid.New()
	tokenString, err := store.SignToken(userID, "user@example.com")
	if err != nil {
		t.Fatalf("failed to sign token: %v", err)
	}

	router := newTestRouter()
	router.Use(AuthMiddleware())

	var gotUserID uuid.UUID
	var gotToken string
	router.GET("/protected", func(c *gin.Context) {
		gotUserID = c.MustGet("user_id").(uuid.UUID)
		gotToken = c.MustGet("token").(string)
		c.String(http.StatusOK, "ok")
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
	if gotUserID != userID {
		t.Errorf("expected user_id %s in context, got %s", userID, gotUserID)
	}
	if gotToken != tokenString {
		t.Errorf("expected token %q in context, got %q", tokenString, gotToken)
	}
}

func TestLogResponse_PassesThroughAndCapturesBody(t *testing.T) {
	router := newTestRouter()
	router.Use(LogResponse())
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusCreated, gin.H{"message": "hello"})
	})

	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", w.Code)
	}
	if got := w.Body.String(); got == "" {
		t.Error("expected response body to be forwarded to the client")
	}
}
