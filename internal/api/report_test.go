package api

import "testing"

func TestValidUTF8(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"ascii", "hello", "hello"},
		{"valid utf8", "你好 world", "你好 world"},
		{"empty", "", ""},
		{"invalid run collapsed", "\xc4\xe3\xba\xc3", "\uFFFD"},
		{"invalid suffix", "ok\xe8", "ok\uFFFD"},
		{"dangling lead before replacement char", "\xe8\xef\xbf\xbd", "\uFFFD\uFFFD"},
	}
	for _, tt := range tests {
		if got := validUTF8(tt.in); got != tt.want {
			t.Errorf("validUTF8(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}
