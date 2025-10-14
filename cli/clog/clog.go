// Customise log for the print
package clog

import (
	"fmt"

	"github.com/fatih/color"
)

// Log prints the message without styling (like println!("{}", message)).

func Log(message string) {
	fmt.Fprintln(color.Output, message)
}

// ILog prints a bold cyan ">" followed by the message.
func ILog(message string) {
	mark := color.New(color.FgCyan, color.Bold).Sprint(">")
	fmt.Fprintf(color.Output, "%s %s\n", mark, message)
}

// WLog prints a bold orange "!" followed by the message.
func WLog(message string) {
	mark := color.New(color.FgHiMagenta, color.Bold).Sprint("!")
	fmt.Fprintf(color.Output, "%s %s\n", mark, message)
}

// ELog prints a bold red "🞪" followed by the message.
func ELog(message string) {
	mark := color.New(color.FgRed, color.Bold).Sprint("🞪")
	fmt.Fprintf(color.Output, "%s %s\n", mark, message)
}

// SLog prints a bold green "✓" followed by the message.
func SLog(message string) {
	mark := color.New(color.FgGreen, color.Bold).Sprint("✓")
	fmt.Fprintf(color.Output, "%s %s\n", mark, message)
}
