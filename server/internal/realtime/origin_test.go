package realtime

import (
	"net/http"
	"reflect"
	"testing"
)

func TestNormalizeOrigin(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{name: "trim and lowercase", in: " HTTPS://APP.EXAMPLE.COM ", want: "https://app.example.com"},
		{name: "drop default https port", in: "https://app.example.com:443", want: "https://app.example.com"},
		{name: "drop default http port", in: "http://app.example.com:80", want: "http://app.example.com"},
		{name: "keep custom port", in: "http://app.example.com:3000", want: "http://app.example.com:3000"},
		{name: "invalid stays trimmed", in: "not-a-url", want: "not-a-url"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := normalizeOrigin(tt.in); got != tt.want {
				t.Fatalf("normalizeOrigin() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestResolveAllowedOriginsFromEnvPriority(t *testing.T) {
	t.Setenv("ALLOWED_ORIGINS", "https://ws.example.com")
	t.Setenv("CORS_ALLOWED_ORIGINS", "https://cors.example.com")
	t.Setenv("FRONTEND_ORIGIN", "https://front.example.com")
	got := ResolveAllowedOriginsFromEnv()
	want := []string{"https://ws.example.com"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("ResolveAllowedOriginsFromEnv() = %#v, want %#v", got, want)
	}
}

func TestResolveAllowedOriginsFromEnvFallbacks(t *testing.T) {
	t.Run("from cors", func(t *testing.T) {
		t.Setenv("ALLOWED_ORIGINS", "")
		t.Setenv("CORS_ALLOWED_ORIGINS", "https://cors.example.com, https://api.example.com:443")
		t.Setenv("FRONTEND_ORIGIN", "")
		got := ResolveAllowedOriginsFromEnv()
		want := []string{"https://cors.example.com", "https://api.example.com"}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("ResolveAllowedOriginsFromEnv() = %#v, want %#v", got, want)
		}
	})

	t.Run("from frontend", func(t *testing.T) {
		t.Setenv("ALLOWED_ORIGINS", "")
		t.Setenv("CORS_ALLOWED_ORIGINS", "")
		t.Setenv("FRONTEND_ORIGIN", "https://front.example.com")
		got := ResolveAllowedOriginsFromEnv()
		want := []string{"https://front.example.com"}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("ResolveAllowedOriginsFromEnv() = %#v, want %#v", got, want)
		}
	})

	t.Run("defaults", func(t *testing.T) {
		t.Setenv("ALLOWED_ORIGINS", "")
		t.Setenv("CORS_ALLOWED_ORIGINS", "")
		t.Setenv("FRONTEND_ORIGIN", "")
		got := ResolveAllowedOriginsFromEnv()
		want := []string{
			"http://localhost:3000",
			"http://localhost:5173",
			"http://localhost:5174",
		}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("ResolveAllowedOriginsFromEnv() = %#v, want %#v", got, want)
		}
	})
}

func TestCheckOriginNormalizedMatch(t *testing.T) {
	previous := allowedWSOrigins.Load().([]string)
	t.Cleanup(func() { allowedWSOrigins.Store(previous) })
	allowedWSOrigins.Store([]string{"https://app.example.com"})

	req, err := http.NewRequest("GET", "http://localhost/ws", nil)
	if err != nil {
		t.Fatalf("NewRequest: %v", err)
	}
	req.Header.Set("Origin", "https://APP.EXAMPLE.COM:443")
	if !checkOrigin(req) {
		t.Fatal("checkOrigin() = false, want true for normalized match")
	}
}

