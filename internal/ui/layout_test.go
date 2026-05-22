package ui

import "testing"

func TestMainLayout_wideUsesThreeColumns(t *testing.T) {
	l := MainLayoutFor(120, 35)
	if l.Mode != LayoutModeWide || l.OutpostBelow {
		t.Fatalf("got mode=%v outpostBelow=%v", l.Mode, l.OutpostBelow)
	}
}

func TestMainLayout_logWidthOnWide(t *testing.T) {
	if w := LogViewportWidth(120, 35); w < 44 {
		t.Fatalf("log width=%d on wide terminal", w)
	}
}
