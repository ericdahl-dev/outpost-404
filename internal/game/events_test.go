package game

import "testing"

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
	if len(s.Log) == 0 || s.Log[0] != "Quiet Shift: Calm day." {
		t.Fatalf("log[0] = %q, want standard event line", s.Log[0])
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

func TestTriggerRandomEvent_LiveMatchesReplayByEventID(t *testing.T) {
	content := testContentWithEvents()
	live := NewStateWithSeed(content, 42)
	live.NextDay()
	if len(live.Log) == 0 {
		t.Fatal("expected an event log line from NextDay")
	}
	var eventID string
	for _, e := range content.Events {
		if live.Log[0] == e.Title+": "+e.Description {
			eventID = e.ID
			break
		}
	}
	if eventID == "" {
		t.Fatalf("could not match log %q to a test event", live.Log[0])
	}

	replay := NewStateWithSeed(content, 42)
	replay.advanceDay()
	replay.applyEventByID(eventID)
	replay.Clamp()

	if diff := snapshotDiff(normalizeSnapshot(live.snapshot()), normalizeSnapshot(replay.snapshot())); diff != "" {
		t.Fatalf("live NextDay vs replay event_id: %s", diff)
	}
}
