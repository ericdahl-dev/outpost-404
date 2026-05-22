package ui

import (
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

func TestFormatColonyLogLine_DangerPrefix(t *testing.T) {
	t.Setenv("NO_COLOR", "")
	prev := lipgloss.ColorProfile()
	t.Cleanup(func() { lipgloss.SetColorProfile(prev) })
	lipgloss.SetColorProfile(termenv.TrueColor)

	got := FormatColonyLogLine("! Power reserves critical")
	if got == "! Power reserves critical" || !strings.Contains(got, "! Power reserves critical") {
		t.Fatalf("expected styled danger line, got %q", got)
	}
}

func TestFormatColonyLogLine_EventPrefix(t *testing.T) {
	t.Setenv("NO_COLOR", "")
	prev := lipgloss.ColorProfile()
	t.Cleanup(func() { lipgloss.SetColorProfile(prev) })
	lipgloss.SetColorProfile(termenv.TrueColor)

	got := FormatColonyLogLine("> Solar Storm: flare. Power -9.")
	if got == "> Solar Storm: flare. Power -9." || !strings.Contains(got, "> Solar Storm") {
		t.Fatalf("expected styled event line, got %q", got)
	}
}

func TestFormatColonyLogLine_PlainUnchanged(t *testing.T) {
	line := "Survive 45 days or complete 5 beacon parts."
	if FormatColonyLogLine(line) != line {
		t.Fatal("plain system copy should not be restyled")
	}
}

func TestFormatColonyLogLines_JoinsMultiple(t *testing.T) {
	body := FormatColonyLogLines([]string{"+ Built Solar Array level 1.", "> Quiet Shift: calm."})
	if body == "" || len(body) < 20 {
		t.Fatalf("unexpected body: %q", body)
	}
}
