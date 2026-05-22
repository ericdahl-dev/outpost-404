package game

import (
	"strings"
	"testing"
)

func TestFormatEffectSummary_Empty(t *testing.T) {
	if got := formatEffectSummary(nil); got != "" {
		t.Fatalf("got %q, want empty", got)
	}
}

func TestFormatEffectSummary_OrdersResources(t *testing.T) {
	got := formatEffectSummary(map[string]int{
		"credits": 30,
		"power":   -9,
		"morale":  8,
	})
	want := "Power -9, Morale +8, Credits +30"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestFormatEventLogLine_FlavorAndImpact(t *testing.T) {
	line := formatEventLogLine(EventDef{
		Title:       "Solar Storm",
		Description: "Charged particles battered the arrays.",
		Effects:     map[string]int{"power": -9},
	})
	want := "Solar Storm: Charged particles battered the arrays. Power -9."
	if line != want {
		t.Fatalf("got %q, want %q", line, want)
	}
}

func TestApplyEvent_LogIncludesFlavorAndEffectSummary(t *testing.T) {
	s := newTestState()
	s.applyEvent(EventDef{
		Title:       "Quiet Shift",
		Description: "Nothing broke. People noticed.",
		Effects:     map[string]int{"morale": 8, "food": 2},
	})
	if len(s.Log) == 0 {
		t.Fatal("expected log line")
	}
	line := s.Log[0]
	if !logContainsAll(line, "Quiet Shift", "Nothing broke", "Morale +8", "Food +2") {
		t.Fatalf("log = %q, missing flavor or impact", line)
	}
}

func TestCheckEnd_CollapseLogsCause(t *testing.T) {
	s := newTestState()
	s.Food = 0
	s.CheckEnd()
	if !s.GameOver || s.Won {
		t.Fatal("expected collapse")
	}
	if len(s.Log) == 0 || !logContainsAll(s.Log[0], "*", "food reserves depleted") {
		t.Fatalf("log[0] = %q, want milestone collapse cause", s.Log[0])
	}
}

func TestSessionSummary_IncludesKeyMoments(t *testing.T) {
	s := newTestState()
	s.AddLogKind(LogMilestone, "Beacon part 1/5 installed.")
	s.AddLogKind(LogMilestone, "Hydroponics came online.")
	summary := s.SessionSummary()
	if len(summary) < 2 {
		t.Fatalf("summary too short: %v", summary)
	}
	joined := joinStrings(summary)
	if !logContainsAll(joined, "Beacon part 1/5", "Hydroponics") {
		t.Fatalf("summary missing key moments: %v", summary)
	}
}

func logContainsAll(s string, parts ...string) bool {
	for _, p := range parts {
		if !strings.Contains(s, p) {
			return false
		}
	}
	return true
}

func joinStrings(parts []string) string {
	return strings.Join(parts, "\n")
}
