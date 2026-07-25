package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// newTestHandler returns a Handler with a nil store. This is safe for tests
// that exercise validation failures returned before any store access occurs.
func newTestHandler() *Handler {
	return NewHandler(nil)
}

func doJSONRequest(router *gin.Engine, method, path string, body interface{}) *httptest.ResponseRecorder {
	var reqBody *bytes.Buffer
	if body != nil {
		b, _ := json.Marshal(body)
		reqBody = bytes.NewBuffer(b)
	} else {
		reqBody = bytes.NewBuffer(nil)
	}
	req := httptest.NewRequest(method, path, reqBody)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

func decodeErrorResponse(t *testing.T, w *httptest.ResponseRecorder) map[string]interface{} {
	t.Helper()
	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response body: %v (body=%s)", err, w.Body.String())
	}
	return resp
}

func TestRegister_InvalidJSON(t *testing.T) {
	h := newTestHandler()
	router := gin.New()
	router.POST("/register", h.Register)

	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBufferString("not-json"))
	req.Header.Set("Content-Type", "application/json")
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

func TestRegister_MissingEmail(t *testing.T) {
	h := newTestHandler()
	router := gin.New()
	router.POST("/register", h.Register)

	w := doJSONRequest(router, http.MethodPost, "/register", map[string]string{
		"password": "password123",
	})

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d, body=%s", w.Code, w.Body.String())
	}
}

func TestRegister_InvalidEmailFormat(t *testing.T) {
	h := newTestHandler()
	router := gin.New()
	router.POST("/register", h.Register)

	w := doJSONRequest(router, http.MethodPost, "/register", map[string]string{
		"email":    "not-an-email",
		"password": "password123",
	})

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d, body=%s", w.Code, w.Body.String())
	}
}

func TestRegister_PasswordTooShort(t *testing.T) {
	h := newTestHandler()
	router := gin.New()
	router.POST("/register", h.Register)

	w := doJSONRequest(router, http.MethodPost, "/register", map[string]string{
		"email":    "user@example.com",
		"password": "abc",
	})

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d, body=%s", w.Code, w.Body.String())
	}
}

func TestLogin_InvalidJSON(t *testing.T) {
	h := newTestHandler()
	router := gin.New()
	router.POST("/login", h.Login)

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBufferString("{"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d, body=%s", w.Code, w.Body.String())
	}
}

func TestLogin_MissingPassword(t *testing.T) {
	h := newTestHandler()
	router := gin.New()
	router.POST("/login", h.Login)

	w := doJSONRequest(router, http.MethodPost, "/login", map[string]string{
		"email": "user@example.com",
	})

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d, body=%s", w.Code, w.Body.String())
	}
}

func TestLogin_InvalidEmailFormat(t *testing.T) {
	h := newTestHandler()
	router := gin.New()
	router.POST("/login", h.Login)

	w := doJSONRequest(router, http.MethodPost, "/login", map[string]string{
		"email":    "not-an-email",
		"password": "whatever",
	})

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d, body=%s", w.Code, w.Body.String())
	}
}

func TestSaveKundali_MissingUserContext(t *testing.T) {
	h := newTestHandler()
	router := gin.New()
	// Intentionally not setting user_id in context, simulating a request
	// that bypassed AuthMiddleware.
	router.POST("/save-kundali", h.SaveKundali)

	w := doJSONRequest(router, http.MethodPost, "/save-kundali", map[string]string{})

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d, body=%s", w.Code, w.Body.String())
	}
	resp := decodeErrorResponse(t, w)
	if resp["error"] != "unauthorized" {
		t.Errorf("expected error 'unauthorized', got %v", resp["error"])
	}
}

// withUserContext stubs an authenticated request by injecting a user_id into
// the gin context, simulating what AuthMiddleware would normally set.
func withUserContext(h gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("user_id", uuid.New())
		h(c)
	}
}

func TestSaveKundali_MissingRequiredField(t *testing.T) {
	h := newTestHandler()
	router := gin.New()
	router.POST("/save-kundali", withUserContext(h.SaveKundali))

	// "day" is omitted, which should fail the model.Kundli binding
	// validation before any store access occurs.
	w := doJSONRequest(router, http.MethodPost, "/save-kundali", map[string]string{
		"name":       "Arjun",
		"month":      "08",
		"year":       "1995",
		"hour":       "10",
		"minute":     "30",
		"birthPlace": "Delhi, India",
		"latitude":   "28.6139",
		"longitude":  "77.2090",
		"timeZone":   "+05:30",
		"isFemale":   "Male",
	})

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d, body=%s", w.Code, w.Body.String())
	}
	resp := decodeErrorResponse(t, w)
	if resp["error"] != "validation_error" {
		t.Errorf("expected error 'validation_error', got %v", resp["error"])
	}
}

func TestSaveKundali_EmptyBody(t *testing.T) {
	h := newTestHandler()
	router := gin.New()
	router.POST("/save-kundali", withUserContext(h.SaveKundali))

	w := doJSONRequest(router, http.MethodPost, "/save-kundali", map[string]string{})

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d, body=%s", w.Code, w.Body.String())
	}
}
