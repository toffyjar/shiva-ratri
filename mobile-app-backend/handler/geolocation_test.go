package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestSearchPlaces_MissingQuery(t *testing.T) {
	h := newTestHandler()
	router := gin.New()
	router.GET("/geolocation/search", h.SearchPlaces)

	req := httptest.NewRequest(http.MethodGet, "/geolocation/search", nil)
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

func TestSearchPlaces_EmptyQueryParam(t *testing.T) {
	h := newTestHandler()
	router := gin.New()
	router.GET("/geolocation/search", h.SearchPlaces)

	req := httptest.NewRequest(http.MethodGet, "/geolocation/search?q=", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d, body=%s", w.Code, w.Body.String())
	}
}
