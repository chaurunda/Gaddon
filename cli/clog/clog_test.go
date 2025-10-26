package clog

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/fatih/color"
)

// captureOutput is a helper function to capture output from color.Output
func captureOutput(f func()) string {
	// Save original output
	originalOutput := color.Output

	// Create a buffer to capture output
	var buf bytes.Buffer
	color.Output = &buf

	// Execute the function
	f()

	// Restore original output
	color.Output = originalOutput

	return buf.String()
}

func TestLog(t *testing.T) {
	testMessage := "This is a test message"

	output := captureOutput(func() {
		Log(testMessage)
	})

	// Log should print the message followed by a newline
	expected := testMessage + "\n"
	if output != expected {
		t.Errorf("Log() output = %q, expected %q", output, expected)
	}
}

func TestILog(t *testing.T) {
	testMessage := "Info message"

	output := captureOutput(func() {
		ILog(testMessage)
	})

	// ILog should contain the message and a ">" marker
	if !strings.Contains(output, testMessage) {
		t.Errorf("ILog() output should contain message %q, got %q", testMessage, output)
	}

	// Should contain the ">" marker (without color codes for this test)
	if !strings.Contains(output, ">") {
		t.Errorf("ILog() output should contain '>' marker, got %q", output)
	}

	// Should end with newline
	if !strings.HasSuffix(output, "\n") {
		t.Errorf("ILog() output should end with newline, got %q", output)
	}
}

func TestWLog(t *testing.T) {
	testMessage := "Warning message"

	output := captureOutput(func() {
		WLog(testMessage)
	})

	// WLog should contain the message and a "!" marker
	if !strings.Contains(output, testMessage) {
		t.Errorf("WLog() output should contain message %q, got %q", testMessage, output)
	}

	// Should contain the "!" marker
	if !strings.Contains(output, "!") {
		t.Errorf("WLog() output should contain '!' marker, got %q", output)
	}

	// Should end with newline
	if !strings.HasSuffix(output, "\n") {
		t.Errorf("WLog() output should end with newline, got %q", output)
	}
}

func TestELog(t *testing.T) {
	testMessage := "Error message"

	output := captureOutput(func() {
		ELog(testMessage)
	})

	// ELog should contain the message and a "🞪" marker
	if !strings.Contains(output, testMessage) {
		t.Errorf("ELog() output should contain message %q, got %q", testMessage, output)
	}

	// Should contain the "🞪" marker
	if !strings.Contains(output, "🞪") {
		t.Errorf("ELog() output should contain '🞪' marker, got %q", output)
	}

	// Should end with newline
	if !strings.HasSuffix(output, "\n") {
		t.Errorf("ELog() output should end with newline, got %q", output)
	}
}

func TestSLog(t *testing.T) {
	testMessage := "Success message"

	output := captureOutput(func() {
		SLog(testMessage)
	})

	// SLog should contain the message and a "✓" marker
	if !strings.Contains(output, testMessage) {
		t.Errorf("SLog() output should contain message %q, got %q", testMessage, output)
	}

	// Should contain the "✓" marker
	if !strings.Contains(output, "✓") {
		t.Errorf("SLog() output should contain '✓' marker, got %q", output)
	}

	// Should end with newline
	if !strings.HasSuffix(output, "\n") {
		t.Errorf("SLog() output should end with newline, got %q", output)
	}
}

func TestLogFunctionsWithEmptyMessage(t *testing.T) {
	tests := []struct {
		name     string
		function func(string)
		marker   string
	}{
		{"Log", Log, ""},
		{"ILog", ILog, ">"},
		{"WLog", WLog, "!"},
		{"ELog", ELog, "🞪"},
		{"SLog", SLog, "✓"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := captureOutput(func() {
				tt.function("")
			})

			if tt.marker != "" {
				// For functions with markers, should still contain the marker
				if !strings.Contains(output, tt.marker) {
					t.Errorf("%s() with empty message should still contain marker %q, got %q", tt.name, tt.marker, output)
				}
			}

			// All functions should end with newline
			if !strings.HasSuffix(output, "\n") {
				t.Errorf("%s() output should end with newline, got %q", tt.name, output)
			}
		})
	}
}

func TestLogFunctionsWithSpecialCharacters(t *testing.T) {
	specialMessage := "Message with special chars: !@#$%^&*()[]{}|\\:;\"'<>,.?/~`"

	tests := []struct {
		name     string
		function func(string)
		marker   string
	}{
		{"Log", Log, ""},
		{"ILog", ILog, ">"},
		{"WLog", WLog, "!"},
		{"ELog", ELog, "🞪"},
		{"SLog", SLog, "✓"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := captureOutput(func() {
				tt.function(specialMessage)
			})

			// Should contain the special message
			if !strings.Contains(output, specialMessage) {
				t.Errorf("%s() output should contain special message, got %q", tt.name, output)
			}

			if tt.marker != "" {
				// Should contain the marker
				if !strings.Contains(output, tt.marker) {
					t.Errorf("%s() output should contain marker %q, got %q", tt.name, tt.marker, output)
				}
			}
		})
	}
}

func TestLogFunctionsWithMultilineMessage(t *testing.T) {
	multilineMessage := "Line 1\nLine 2\nLine 3"

	tests := []struct {
		name     string
		function func(string)
		marker   string
	}{
		{"Log", Log, ""},
		{"ILog", ILog, ">"},
		{"WLog", WLog, "!"},
		{"ELog", ELog, "🞪"},
		{"SLog", SLog, "✓"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := captureOutput(func() {
				tt.function(multilineMessage)
			})

			// Should contain all lines of the message
			lines := strings.Split(multilineMessage, "\n")
			for _, line := range lines {
				if !strings.Contains(output, line) {
					t.Errorf("%s() output should contain line %q, got %q", tt.name, line, output)
				}
			}

			if tt.marker != "" {
				// Should contain the marker
				if !strings.Contains(output, tt.marker) {
					t.Errorf("%s() output should contain marker %q, got %q", tt.name, tt.marker, output)
				}
			}
		})
	}
}

// Benchmark tests to ensure performance is reasonable
func BenchmarkLog(b *testing.B) {
	// Redirect to discard to avoid console output during benchmarks
	color.Output = io.Discard
	defer func() { color.Output = os.Stdout }()

	message := "Benchmark test message"
	for i := 0; i < b.N; i++ {
		Log(message)
	}
}

func BenchmarkILog(b *testing.B) {
	color.Output = io.Discard
	defer func() { color.Output = os.Stdout }()

	message := "Benchmark test message"
	for i := 0; i < b.N; i++ {
		ILog(message)
	}
}

func BenchmarkWLog(b *testing.B) {
	color.Output = io.Discard
	defer func() { color.Output = os.Stdout }()

	message := "Benchmark test message"
	for i := 0; i < b.N; i++ {
		WLog(message)
	}
}

func BenchmarkELog(b *testing.B) {
	color.Output = io.Discard
	defer func() { color.Output = os.Stdout }()

	message := "Benchmark test message"
	for i := 0; i < b.N; i++ {
		ELog(message)
	}
}

func BenchmarkSLog(b *testing.B) {
	color.Output = io.Discard
	defer func() { color.Output = os.Stdout }()

	message := "Benchmark test message"
	for i := 0; i < b.N; i++ {
		SLog(message)
	}
}