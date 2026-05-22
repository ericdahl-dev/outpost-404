package ui_test

import (
	"testing"

	"github.com/ericdahl/outpost-404/internal/ui"
)

func TestMainLayout_wideTerminalUsesHorizontalPanels(t *testing.T) {
	l := ui.MainLayoutFor(120, 30)
	if l.Stacked {
		t.Fatal("expected horizontal layout at 120 columns")
	}
	if l.LeftWidth < 28 || l.MiddleWidth < 28 {
		t.Fatalf("panel widths too small: left=%d mid=%d", l.LeftWidth, l.MiddleWidth)
	}
	if l.LogInnerWidth < 44 {
		t.Fatalf("log width=%d, want at least 44 on wide terminal", l.LogInnerWidth)
	}
}

func TestMainLayout_narrowTerminalStacksPanels(t *testing.T) {
	l := ui.MainLayoutFor(70, 24)
	if !l.Stacked {
		t.Fatal("expected stacked layout below horizontal breakpoint")
	}
	if l.ContentWidth < 60 {
		t.Fatalf("content width=%d, want most of terminal width", l.ContentWidth)
	}
}

func TestMainLayout_logViewportShrinksOnNarrowTerminal(t *testing.T) {
	w := ui.LogViewportWidth(60)
	if w > 52 {
		t.Fatalf("log viewport width=%d, should shrink with terminal", w)
	}
	if w < 20 {
		t.Fatalf("log viewport width=%d, too narrow to read", w)
	}
}
