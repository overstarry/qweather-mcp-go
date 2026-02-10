package api

import "testing"

func TestLogLevelString(t *testing.T) {
	tests := []struct {
		level LogLevel
		want  string
	}{
		{LogLevelNone, "NONE"},
		{LogLevelError, "ERROR"},
		{LogLevelInfo, "INFO"},
		{LogLevelDebug, "DEBUG"},
		{LogLevel(999), "UNKNOWN"},
	}

	for _, tt := range tests {
		if got := tt.level.String(); got != tt.want {
			t.Fatalf("LogLevel(%d).String() = %q, want %q", tt.level, got, tt.want)
		}
	}
}
