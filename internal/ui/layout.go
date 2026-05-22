package ui

// LayoutMode selects the main-screen panel arrangement.
type LayoutMode int

const (
	LayoutModeNarrow LayoutMode = iota
	LayoutModeMedium
	LayoutModeWide
)

// MainLayout describes how the main screen arranges panels for the current terminal size.
type MainLayout struct {
	Mode          LayoutMode
	Stacked       bool
	OutpostBelow  bool
	ContentWidth  int
	LeftWidth     int
	MiddleWidth   int
	LogInnerWidth int
	BoxWidth      int
	CompactKeys   bool
}

const (
	wideMinWidth    = 120
	wideMinHeight   = 35
	mediumMinWidth  = 100
	mediumMinHeight = 30
	minPanelWidth   = 24
	maxLeftWidth    = 36
	maxMiddleWidth  = 42
)

func MainLayoutFor(termWidth, termHeight int) MainLayout {
	content := contentWidth(termWidth)
	logW := LogViewportWidth(termWidth, termHeight)
	l := MainLayout{
		ContentWidth:  content,
		LogInnerWidth: logW,
		BoxWidth:      boxWidth(termWidth, 80),
		CompactKeys:   termWidth < 88,
	}
	switch {
	case termWidth >= wideMinWidth && termHeight >= wideMinHeight:
		l.Mode = LayoutModeWide
		l.applyWideLayout(termWidth, logW)
	case termWidth >= mediumMinWidth && termHeight >= mediumMinHeight:
		l.Mode = LayoutModeMedium
		l.applyMediumLayout(termWidth, logW, content)
	default:
		l.Mode = LayoutModeNarrow
		l.applyNarrowLayout(content)
	}
	return l
}

func (l *MainLayout) applyWideLayout(termWidth, logW int) {
	l.Stacked = false
	l.OutpostBelow = false
	logBox := logW + 2
	remaining := termWidth - logBox
	if remaining < minPanelWidth*2 {
		l.applyNarrowLayout(l.ContentWidth)
		l.Mode = LayoutModeNarrow
		return
	}
	l.LeftWidth = clamp(minPanelWidth, remaining*9/20, maxLeftWidth)
	l.MiddleWidth = clamp(minPanelWidth, remaining-l.LeftWidth, maxMiddleWidth)
}

func (l *MainLayout) applyMediumLayout(termWidth, logW, content int) {
	l.Stacked = false
	l.OutpostBelow = true
	l.MiddleWidth = content
	logBox := logW + 2
	remaining := termWidth - logBox
	if remaining < minPanelWidth {
		l.applyNarrowLayout(content)
		l.Mode = LayoutModeNarrow
		return
	}
	l.LeftWidth = clamp(minPanelWidth, remaining, maxLeftWidth)
}

func (l *MainLayout) applyNarrowLayout(content int) {
	l.Stacked = true
	l.OutpostBelow = false
	l.LeftWidth = content
	l.MiddleWidth = content
}

func contentWidth(termWidth int) int {
	w := termWidth - 4
	if w < 40 {
		return 40
	}
	return w
}

func boxWidth(termWidth, preferred int) int {
	w := termWidth - 4
	if w < 40 {
		return 40
	}
	if w > preferred {
		return preferred
	}
	return w
}

func LogViewportWidth(termWidth, termHeight int) int {
	if termWidth < mediumMinWidth || termHeight < mediumMinHeight {
		return clamp(22, termWidth-8, 44)
	}
	w := 44
	if termWidth >= wideMinWidth {
		return 52
	}
	return w
}

func LogViewportHeight(termHeight int) int {
	h := termHeight - 14
	return clamp(6, h, 14)
}

func clamp(min, v, max int) int {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}
