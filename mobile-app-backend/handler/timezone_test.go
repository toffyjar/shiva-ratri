package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestSearchTimezone_MissingLatitudeAndLongitude(t *testing.T) {
	h := newTestHandler()
	router := gin.New()
	router.GET("/timezone/search", h.SearchTimezone)

	req := httptest.NewRequest(http.MethodGet, "/timezone/search", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d, body=%s", w.Code, w.Body.String())
	}
	resp := decodeErrorResponse(t, w)
	if resp["error"] != "validation_error" {
		t.Errorf("expected error 'validation_error', got %v", resp["error"])
	}
}

func TestSearchTimezone_MissingLongitude(t *testing.T) {
	h := newTestHandler()
	router := gin.New()
	router.GET("/timezone/search", h.SearchTimezone)

	req := httptest.NewRequest(http.MethodGet, "/timezone/search?latitude=12.34", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d, body=%s", w.Code, w.Body.String())
	}
}

func TestSearchTimezone_MissingLatitude(t *testing.T) {
	h := newTestHandler()
	router := gin.New()
	router.GET("/timezone/search", h.SearchTimezone)

	req := httptest.NewRequest(http.MethodGet, "/timezone/search?longitude=56.78", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d, body=%s", w.Code, w.Body.String())
	}
}

func TestSearchTimezone_NonNumericLatitude(t *testing.T) {
	h := newTestHandler()
	router := gin.New()
	router.GET("/timezone/search", h.SearchTimezone)

	req := httptest.NewRequest(http.MethodGet, "/timezone/search?latitude=abc&longitude=56.78", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d, body=%s", w.Code, w.Body.String())
	}
}

func TestSearchTimezone_NonNumericLongitude(t *testing.T) {
	h := newTestHandler()
	router := gin.New()
	router.GET("/timezone/search", h.SearchTimezone)

	req := httptest.NewRequest(http.MethodGet, "/timezone/search?latitude=12.34&longitude=xyz", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d, body=%s", w.Code, w.Body.String())
	}
}
