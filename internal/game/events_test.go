package game

import (
	"strings"
	"testing"
)

func TestApplyEvent_AppliesEffectsAndLogLine(t *testing.T) {
	s := newTestState()
	startMorale := s.Morale
	event := EventDef{
		ID:          "quiet_shift",
		Title:       "Quiet Shift",
		Description: "Calm day.",
		Effects:     map[string]int{"morale": 8},
		MinDay:      1,
	}

	s.applyEvent(event)

	if s.Morale != startMorale+8 {
		t.Fatalf("Morale = %d, want %d", s.Morale, startMorale+8)
	}
	if len(s.Log) == 0 || !strings.Contains(s.Log[0], "Quiet Shift: Calm day.") || !strings.Contains(s.Log[0], "Morale +8") {
		t.Fatalf("log[0] = %q, want flavor and effect summary", s.Log[0])
	}
}

func TestApplyEventByID_MatchesApplyEvent(t *testing.T) {
	content := testContentWithEvents()
	event := content.Events[0]

	byID := NewState(content)
	byID.Power = 50
	byID.Food = 50
	byID.applyEventByID(event.ID)

	direct := NewState(content)
	direct.Power = 50
	direct.Food = 50
	direct.applyEvent(event)

	if diff := snapshotDiff(byID.snapshot(), direct.snapshot()); diff != "" {
		t.Fatalf("applyEventByID vs applyEvent: %s", diff)
	}
}

func TestApplyEventByID_UnknownID_IsNoOp(t *testing.T) {
	s := newTestState()
	before := s.snapshot()
	s.applyEventByID("missing_event")
	if diff := snapshotDiff(before, s.snapshot()); diff != "" {
		t.Fatalf("unknown event changed state: %s", diff)
	}
}

func TestEmbeddedContent_DustSealFailureEvent(t *testing.T) {
	content, err := LoadEmbeddedContent()
	if err != nil {
		t.Fatalf("LoadEmbeddedContent: %v", err)
	}
	var dust *EventDef
	for i := range content.Events {
		if content.Events[i].ID == "dust_seal_failure" {
			dust = &content.Events[i]
			break
		}
	}
	if dust == nil {
		t.Fatal("expected dust_seal_failure in embedded events")
	}
	if dust.MinDay != 14 {
		t.Fatalf("MinDay = %d, want 14", dust.MinDay)
	}
	if dust.Effects["food"] != -6 || dust.Effects["morale"] != -3 {
		t.Fatalf("effects = %v, want food -6 morale -3", dust.Effects)
	}
}

func TestApplyEvent_DustSealFailureReducesFoodAndMorale(t *testing.T) {
	s := newTestState()
	startFood, startMorale := s.Food, s.Morale
	event := EventDef{
		ID:          "dust_seal_failure",
		Title:       "Dust Seal Failure",
		Description: "Seals failed.",
		Effects:     map[string]int{"food": -6, "morale": -3},
		MinDay:      14,
	}

	s.applyEvent(event)

	if s.Food != startFood-6 {
		t.Fatalf("Food = %d, want %d", s.Food, startFood-6)
	}
	if s.Morale != startMorale-3 {
		t.Fatalf("Morale = %d, want %d", s.Morale, startMorale-3)
	}
}

func TestEligibleEvents_FiltersByMinDay(t *testing.T) {
	events := []EventDef{
		{ID: "early", MinDay: 1},
		{ID: "mid", MinDay: 5},
		{ID: "late", MinDay: 10},
	}
	got := eligibleEventsForState(State{Day: 5}, events)
	if len(got) != 2 {
		t.Fatalf("eligible count = %d, want 2", len(got))
	}
	if got[0].ID != "early" || got[1].ID != "mid" {
		t.Fatalf("eligible = %v, want early and mid", got)
	}
}

func TestEligibleEvents_EmptyWhenNoneQualify(t *testing.T) {
	events := []EventDef{{ID: "late", MinDay: 10}}
	if len(eligibleEventsForState(State{Day: 1}, events)) != 0 {
		t.Fatal("expected no eligible events on day 1")
	}
}

func TestEligibleEvents_FiltersByMaxDay(t *testing.T) {
	events := []EventDef{
		{ID: "early_only", MinDay: 1, MaxDay: 5},
		{ID: "always", MinDay: 1},
	}
	got := eligibleEventsForState(State{Day: 10}, events)
	if len(got) != 1 || got[0].ID != "always" {
		t.Fatalf("day 10 eligible = %v, want only always", got)
	}
	got5 := eligibleEventsForState(State{Day: 5}, events)
	if len(got5) != 2 {
		t.Fatalf("day 5 eligible count = %d, want 2", len(got5))
	}
}

func TestPickRandomEligibleEvent_WeightedFavorsHigherWeight(t *testing.T) {
	events := []EventDef{
		{ID: "light", Weight: 1},
		{ID: "heavy", Weight: 9},
	}
	s := NewStateWithSeed(testContent(), 99)
	counts := map[string]int{}
	const trials = 8000
	for range trials {
		e, ok := s.pickRandomEligibleEvent(events)
		if !ok {
			t.Fatal("expected pick")
		}
		counts[e.ID]++
	}
	heavyShare := float64(counts["heavy"]) / trials
	if heavyShare < 0.85 {
		t.Fatalf("heavy share = %.3f, want ~0.9 with weights 1:9", heavyShare)
	}
}

func TestPickRandomEligibleEvent_EmptyCandidates(t *testing.T) {
	s := NewStateWithSeed(testContentWithEvents(), 1)
	if _, ok := s.pickRandomEligibleEvent(nil); ok {
		t.Fatal("expected false for empty candidates")
	}
}

func TestTriggerRandomEvent_NoEligibleEventsNeverApplies(t *testing.T) {
	for _, seed := range []int64{0, 1, 7, 42, 99} {
		s := NewStateWithSeed(testContent(), seed)
		s.Day = 5
		beforeLog := len(s.Log)
		beforeMorale := s.Morale
		if id := s.TriggerRandomEvent(); id != "" {
			t.Fatalf("seed %d: id=%q with no events in content", seed, id)
		}
		if len(s.Log) != beforeLog {
			t.Fatalf("seed %d: unexpected log line without events", seed)
		}
		if s.Morale != beforeMorale {
			t.Fatalf("seed %d: morale changed without event", seed)
		}
	}
}

func TestTriggerRandomEvent_LiveMatchesReplayByEventID(t *testing.T) {
	content := testContentWithEvents()
	live := NewStateWithSeed(content, 42)
	live.NextDay()
	if len(live.Log) == 0 {
		t.Fatal("expected an event log line from NextDay")
	}
	var eventID string
	for _, e := range content.Events {
		if live.Log[0] == formatEventLogLine(e) {
			eventID = e.ID
			break
		}
	}
	if eventID == "" {
		t.Fatalf("could not match log %q to a test event", live.Log[0])
	}

	replay := NewStateWithSeed(content, 42)
	replay.nextDayWithDetail(map[string]any{"event_id": eventID})

	if diff := snapshotDiff(normalizeSnapshot(live.snapshot()), normalizeSnapshot(replay.snapshot())); diff != "" {
		t.Fatalf("live NextDay vs replay event_id: %s", diff)
	}
}
