package ui

// MainLayout describes how the main screen arranges panels for the current terminal size.
type MainLayout struct {
	Stacked       bool
	ContentWidth  int
	LeftWidth     int
	MiddleWidth   int
	LogInnerWidth int
	BoxWidth      int
	CompactKeys   bool
}

const (
	horizontalMinWidth = 96
	minPanelWidth      = 24
	maxLeftWidth       = 36
	maxMiddleWidth     = 38
)

func MainLayoutFor(termWidth, termHeight int) MainLayout {
	_ = termHeight
	content := contentWidth(termWidth)
	logW := LogViewportWidth(termWidth)
	l := MainLayout{
		ContentWidth:  content,
		LogInnerWidth: logW,
		BoxWidth:      boxWidth(termWidth, 80),
		CompactKeys:   termWidth < 88,
	}
	if termWidth < horizontalMinWidth {
		l.Stacked = true
		l.LeftWidth = content
		l.MiddleWidth = content
		return l
	}
	l.Stacked = false
	logBox := logW + 2
	remaining := termWidth - logBox
	if remaining < minPanelWidth*2 {
		l.Stacked = true
		l.LeftWidth = content
		l.MiddleWidth = content
		return l
	}
	l.LeftWidth = clamp(minPanelWidth, remaining*9/20, maxLeftWidth)
	l.MiddleWidth = clamp(minPanelWidth, remaining-l.LeftWidth, maxMiddleWidth)
	return l
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

func LogViewportWidth(termWidth int) int {
	if termWidth < horizontalMinWidth {
		return clamp(22, termWidth-8, 44)
	}
	w := 44
	if termWidth > 120 {
		return 52
	}
	return w
}

func LogViewportHeight(termHeight int) int {
	h := termHeight - 12
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
