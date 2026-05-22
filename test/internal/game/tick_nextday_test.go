package game_test

import (
	"github.com/ericdahl/outpost-404/internal/game"
	"testing"
)

func TestNextDayWithDetail_ReplaySkipsRandomRollWhenEventIDPresent(t *testing.T) {
	content := testContentWithEvents()
	live := game.NewStateWithSeed(content, 99)
	liveOutcome := game.NextDayWithDetailForTest(&live, nil)

	replay := game.NewStateWithSeed(content, 99)
	replayOutcome := game.NextDayWithDetailForTest(&replay, map[string]any{"event_id": liveOutcome.EventID})

	if replayOutcome.EventID != liveOutcome.EventID {
		t.Fatalf("eventID = %q, want %q", replayOutcome.EventID, liveOutcome.EventID)
	}
	if diff := game.SnapshotDiffForTest(game.NormalizeSnapshotForTest(game.SnapshotForTest(&live)), game.NormalizeSnapshotForTest(game.SnapshotForTest(&replay))); diff != "" {
		t.Fatalf("live vs replay nextDayWithDetail: %s", diff)
	}
}
