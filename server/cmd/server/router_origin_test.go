package main

import "testing"

func TestOriginListContains(t *testing.T) {
	tests := []struct {
		name  string
		csv   string
		value string
		want  bool
	}{
		{name: "exact match", csv: "https://app.example.com,https://api.example.com", value: "https://app.example.com", want: true},
		{name: "case insensitive", csv: "HTTPS://APP.EXAMPLE.COM", value: "https://app.example.com", want: true},
		{name: "trim spaces", csv: " https://app.example.com ", value: "https://app.example.com", want: true},
		{name: "not found", csv: "https://api.example.com", value: "https://app.example.com", want: false},
		{name: "empty value", csv: "https://api.example.com", value: "", want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := originListContains(tt.csv, tt.value); got != tt.want {
				t.Fatalf("originListContains(%q, %q) = %v, want %v", tt.csv, tt.value, got, tt.want)
			}
		})
	}
}

