package store

import (
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func TestHashPassword(t *testing.T) {
	hashed, err := HashPassword("supersecret")
	if err != nil {
		t.Fatalf("HashPassword returned error: %v", err)
	}
	if hashed == "" {
		t.Fatal("expected non-empty hash")
	}
	if hashed == "supersecret" {
		t.Fatal("expected hash to differ from plaintext password")
	}
}

func TestHashPassword_DifferentHashesForSameInput(t *testing.T) {
	h1, err := HashPassword("samepassword")
	if err != nil {
		t.Fatalf("HashPassword returned error: %v", err)
	}
	h2, err := HashPassword("samepassword")
	if err != nil {
		t.Fatalf("HashPassword returned error: %v", err)
	}
	if h1 == h2 {
		t.Fatal("expected bcrypt to salt hashes differently across calls")
	}
}

func TestComparePassword_Success(t *testing.T) {
	hashed, err := HashPassword("mypassword")
	if err != nil {
		t.Fatalf("HashPassword returned error: %v", err)
	}
	if err := ComparePassword(hashed, "mypassword"); err != nil {
		t.Fatalf("expected passwords to match, got error: %v", err)
	}
}

func TestComparePassword_Mismatch(t *testing.T) {
	hashed, err := HashPassword("mypassword")
	if err != nil {
		t.Fatalf("HashPassword returned error: %v", err)
	}
	if err := ComparePassword(hashed, "wrongpassword"); err == nil {
		t.Fatal("expected error for mismatched password, got nil")
	}
}

func TestComparePassword_InvalidHash(t *testing.T) {
	if err := ComparePassword("not-a-bcrypt-hash", "anything"); err == nil {
		t.Fatal("expected error for malformed hash, got nil")
	}
}

func TestGenerateToken(t *testing.T) {
	tok1, err := GenerateToken()
	if err != nil {
		t.Fatalf("GenerateToken returned error: %v", err)
	}
	if len(tok1) != 32 { // 16 bytes hex-encoded
		t.Fatalf("expected token of length 32, got %d (%q)", len(tok1), tok1)
	}

	tok2, err := GenerateToken()
	if err != nil {
		t.Fatalf("GenerateToken returned error: %v", err)
	}
	if tok1 == tok2 {
		t.Fatal("expected two generated tokens to differ")
	}
}

func TestSignToken_And_VerifyToken_RoundTrip(t *testing.T) {
	userID := uuid.New()
	email := "user@example.com"

	tokenString, err := SignToken(userID, email)
	if err != nil {
		t.Fatalf("SignToken returned error: %v", err)
	}
	if tokenString == "" {
		t.Fatal("expected non-empty token string")
	}

	gotID, err := VerifyToken(tokenString)
	if err != nil {
		t.Fatalf("VerifyToken returned error: %v", err)
	}
	if gotID != userID {
		t.Fatalf("expected user id %s, got %s", userID, gotID)
	}
}

func TestVerifyToken_Malformed(t *testing.T) {
	if _, err := VerifyToken("not-a-jwt-token"); err == nil {
		t.Fatal("expected error for malformed token, got nil")
	}
}

func TestVerifyToken_EmptyString(t *testing.T) {
	if _, err := VerifyToken(""); err == nil {
		t.Fatal("expected error for empty token, got nil")
	}
}

func TestVerifyToken_ExpiredToken(t *testing.T) {
	claims := jwt.MapClaims{
		"user_id": uuid.New().String(),
		"email":   "expired@example.com",
		"exp":     time.Now().Add(-1 * time.Hour).Unix(),
		"iat":     time.Now().Add(-2 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(jwtSecret())
	if err != nil {
		t.Fatalf("failed to sign test token: %v", err)
	}

	if _, err := VerifyToken(signed); err == nil {
		t.Fatal("expected error for expired token, got nil")
	}
}

func TestVerifyToken_WrongSigningMethod(t *testing.T) {
	claims := jwt.MapClaims{
		"user_id": uuid.New().String(),
		"email":   "wrongalg@example.com",
		"exp":     time.Now().Add(1 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodNone, claims)
	signed, err := token.SignedString(jwt.UnsafeAllowNoneSignatureType)
	if err != nil {
		t.Fatalf("failed to sign test token: %v", err)
	}

	if _, err := VerifyToken(signed); err == nil {
		t.Fatal("expected error for non-HMAC signing method, got nil")
	}
}

func TestVerifyToken_InvalidUserIDClaim(t *testing.T) {
	claims := jwt.MapClaims{
		"user_id": "not-a-valid-uuid",
		"exp":     time.Now().Add(1 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(jwtSecret())
	if err != nil {
		t.Fatalf("failed to sign test token: %v", err)
	}

	if _, err := VerifyToken(signed); err == nil {
		t.Fatal("expected error for invalid user_id claim, got nil")
	}
}

func TestVerifyToken_MissingUserIDClaim(t *testing.T) {
	claims := jwt.MapClaims{
		"email": "nouser@example.com",
		"exp":   time.Now().Add(1 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(jwtSecret())
	if err != nil {
		t.Fatalf("failed to sign test token: %v", err)
	}

	if _, err := VerifyToken(signed); err == nil {
		t.Fatal("expected error for missing user_id claim, got nil")
	}
}

func TestVerifyToken_SignedWithDifferentSecret(t *testing.T) {
	claims := jwt.MapClaims{
		"user_id": uuid.New().String(),
		"exp":     time.Now().Add(1 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte("a-different-secret"))
	if err != nil {
		t.Fatalf("failed to sign test token: %v", err)
	}

	if _, err := VerifyToken(signed); err == nil {
		t.Fatal("expected error for token signed with different secret, got nil")
	}
}

func TestJwtSecret_FallsBackWhenUnset(t *testing.T) {
	old, had := os.LookupEnv("JWT_SECRET")
	os.Unsetenv("JWT_SECRET")
	defer func() {
		if had {
			os.Setenv("JWT_SECRET", old)
		}
	}()

	if string(jwtSecret()) != "changeme" {
		t.Fatalf("expected fallback secret 'changeme', got %q", string(jwtSecret()))
	}
}

func TestJwtSecret_UsesEnvWhenSet(t *testing.T) {
	old, had := os.LookupEnv("JWT_SECRET")
	os.Setenv("JWT_SECRET", "my-env-secret")
	defer func() {
		if had {
			os.Setenv("JWT_SECRET", old)
		} else {
			os.Unsetenv("JWT_SECRET")
		}
	}()

	if string(jwtSecret()) != "my-env-secret" {
		t.Fatalf("expected env secret, got %q", string(jwtSecret()))
	}
}

func TestJwtSecret_BlankEnvFallsBack(t *testing.T) {
	old, had := os.LookupEnv("JWT_SECRET")
	os.Setenv("JWT_SECRET", "   ")
	defer func() {
		if had {
			os.Setenv("JWT_SECRET", old)
		} else {
			os.Unsetenv("JWT_SECRET")
		}
	}()

	if string(jwtSecret()) != "changeme" {
		t.Fatalf("expected fallback secret for blank env var, got %q", string(jwtSecret()))
	}
}
