package game

import "testing"

func TestNextDayWithDetail_ReplaySkipsRandomRollWhenEventIDPresent(t *testing.T) {
	content := testContentWithEvents()
	live := NewStateWithSeed(content, 99)
	liveOutcome := live.nextDayWithDetail(nil)

	replay := NewStateWithSeed(content, 99)
	replayOutcome := replay.nextDayWithDetail(map[string]any{"event_id": liveOutcome.eventID})

	if replayOutcome.eventID != liveOutcome.eventID {
		t.Fatalf("eventID = %q, want %q", replayOutcome.eventID, liveOutcome.eventID)
	}
	if diff := snapshotDiff(normalizeSnapshot(live.snapshot()), normalizeSnapshot(replay.snapshot())); diff != "" {
		t.Fatalf("live vs replay nextDayWithDetail: %s", diff)
	}
}
