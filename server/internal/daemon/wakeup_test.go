package daemon

import (
	"io"
	"strings"
	"testing"
)

func TestTaskWakeupURL(t *testing.T) {
	tests := []struct {
		name       string
		baseURL    string
		runtimeIDs []string
		want       string
	}{
		{
			name:       "http base",
			baseURL:    "http://localhost:8080",
			runtimeIDs: []string{"runtime-b", "runtime-a"},
			want:       "ws://localhost:8080/api/daemon/ws?runtime_ids=runtime-a%2Cruntime-b",
		},
		{
			name:       "https base",
			baseURL:    "https://api.example.com",
			runtimeIDs: []string{"runtime-1"},
			want:       "wss://api.example.com/api/daemon/ws?runtime_ids=runtime-1",
		},
		{
			name:       "base path",
			baseURL:    "https://api.example.com/multica",
			runtimeIDs: []string{"runtime-1"},
			want:       "wss://api.example.com/multica/api/daemon/ws?runtime_ids=runtime-1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := taskWakeupURL(tt.baseURL, tt.runtimeIDs)
			if err != nil {
				t.Fatalf("taskWakeupURL: %v", err)
			}
			if got != tt.want {
				t.Fatalf("taskWakeupURL() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestSanitizeWakeupURLForLog(t *testing.T) {
	got := sanitizeWakeupURLForLog("wss://api.example.com/api/daemon/ws?runtime_ids=a,b,c&x=1")
	want := "wss://api.example.com/api/daemon/ws?runtime_count=3&x=1"
	if got != want {
		t.Fatalf("sanitizeWakeupURLForLog() = %q, want %q", got, want)
	}
}

func TestReadHandshakeBodyPreview(t *testing.T) {
	got := readHandshakeBodyPreview(io.NopCloser(strings.NewReader("  bad request  ")))
	if got != "bad request" {
		t.Fatalf("readHandshakeBodyPreview() = %q, want %q", got, "bad request")
	}
}
