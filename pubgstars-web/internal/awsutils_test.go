package internal

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"
)

func makeJWT(payload map[string]string) string {
	header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"HS256","typ":"JWT"}`))
	body, _ := json.Marshal(payload)
	claim := base64.RawURLEncoding.EncodeToString(body)
	return header + "." + claim + ".sig"
}

func TestGetUsernameFromJwtToken(t *testing.T) {
	token := makeJWT(map[string]string{"email": "user@example.com"})
	if got := GetUsernameFromJwtToken(token); got != "user@example.com" {
		t.Errorf("got %q, want %q", got, "user@example.com")
	}
}

func TestGetUsernameFromJwtToken_MissingEmail(t *testing.T) {
	token := makeJWT(map[string]string{"sub": "1234"})
	// datum["email"] is nil when key absent; CovertToString(nil) returns "<nil>"
	got := GetUsernameFromJwtToken(token)
	if got != fmt.Sprintf("%v", nil) {
		t.Errorf("GetUsernameFromJwtToken with no email field: got %q, want %q", got, fmt.Sprintf("%v", nil))
	}
}

func TestGetUsernameFromJwtTokenForAdmin(t *testing.T) {
	token := makeJWT(map[string]string{"cognito:username": "adminuser"})
	if got := GetUsernameFromJwtTokenForAdmin(token); got != "adminuser" {
		t.Errorf("got %q, want %q", got, "adminuser")
	}
}

func TestGetUsernameFromJwtToken_MalformedToken(t *testing.T) {
	// Should return empty string without panicking
	if got := GetUsernameFromJwtToken("notavalidtoken"); got != "" {
		t.Errorf("expected empty string for malformed token, got %q", got)
	}
}

func TestCovertToString(t *testing.T) {
	tests := []struct {
		input interface{}
		want  string
	}{
		{42, "42"},
		{3.14, "3.14"},
		{"hello", "hello"},
		{true, "true"},
	}
	for _, tc := range tests {
		if got := CovertToString(tc.input); got != tc.want {
			t.Errorf("CovertToString(%v) = %q, want %q", tc.input, got, tc.want)
		}
	}
}

func TestConvertMillisToString_Format(t *testing.T) {
	// 2024-06-15 10:30:00 UTC in millis
	millis := int64(1718447400000)
	got := ConvertMillisToString(millis)
	// The function formats as "02.01.2006 15:04" in local time.
	// Verify only that the format matches dd.mm.yyyy hh:mm
	parts := strings.Split(got, " ")
	if len(parts) != 2 {
		t.Fatalf("unexpected format: %q", got)
	}
	if len(parts[0]) != 10 || parts[0][2] != '.' || parts[0][5] != '.' {
		t.Errorf("date part format unexpected: %q", parts[0])
	}
	if len(parts[1]) != 5 || parts[1][2] != ':' {
		t.Errorf("time part format unexpected: %q", parts[1])
	}
}

func TestCurrentTimeMillis(t *testing.T) {
	before := time.Now().UnixNano() / 1_000_000
	got := CurrentTimeMillis()
	after := time.Now().UnixNano() / 1_000_000
	if got < before || got > after {
		t.Errorf("CurrentTimeMillis() = %d, expected between %d and %d", got, before, after)
	}
}

func TestContains(t *testing.T) {
	slice := []string{"a", "b", "c"}
	tests := []struct {
		val  string
		want bool
	}{
		{"a", true},
		{"b", true},
		{"c", true},
		{"d", false},
		{"", false},
	}
	for _, tc := range tests {
		if got := Contains(slice, tc.val); got != tc.want {
			t.Errorf("Contains(%q) = %v, want %v", tc.val, got, tc.want)
		}
	}
}

func TestContains_EmptySlice(t *testing.T) {
	if Contains([]string{}, "a") {
		t.Error("Contains on empty slice should return false")
	}
}

// ensure CovertToString uses %v format and not a panic on nil
func TestCovertToString_Nil(t *testing.T) {
	got := CovertToString(nil)
	if got != fmt.Sprintf("%v", nil) {
		t.Errorf("CovertToString(nil) = %q", got)
	}
}
