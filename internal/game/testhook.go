package game

// Test hooks for external tests in test/internal/game (package game_test).

// NextDayOutcomeForTest mirrors nextDayOutcome for replay/event tests.
type NextDayOutcomeForTest struct {
	EventID          string
	PopulationGrowth bool
}

func AdvanceDayForTest(s *State) { s.advanceDay() }

func ApplyBuildingProductionForTest(s *State) { s.applyBuildingProduction() }

func ApplyEventForTest(s *State, e EventDef) { s.applyEvent(e) }

func ApplyEventByIDForTest(s *State, id string) { s.applyEventByID(id) }

func ApplyDailyUpkeepForTest(s *State) { s.applyDailyUpkeep() }

func ApplyMoraleDriftForTest(s *State) { s.applyMoraleDrift() }

func TryColonistArrivalForTest(s *State) { s.tryColonistArrival() }

func ApplyEffectsForTest(s *State, effects map[string]int, mult int) {
	s.applyEffects(effects, mult)
}

func EligibleEventsForTest(s State, events []EventDef) []EventDef {
	return eligibleEventsForState(s, events)
}

func PickRandomEligibleEventForTest(s *State, candidates []EventDef) (EventDef, bool) {
	return s.pickRandomEligibleEvent(candidates)
}

func NextDayWithDetailForTest(s *State, detail map[string]any) NextDayOutcomeForTest {
	o := s.nextDayWithDetail(detail)
	return NextDayOutcomeForTest{EventID: o.eventID, PopulationGrowth: o.populationGrowth}
}

func SnapshotForTest(s *State) Snapshot { return s.snapshot() }

func SeedFromDetailForTest(detail map[string]any) (int64, error) { return seedFromDetail(detail) }

func SnapshotVitalsDepletedForTest(s Snapshot) bool { return snapshotVitalsDepleted(s) }

func SnapshotMeetsSurvivalEndMarginsForTest(s Snapshot) bool {
	return snapshotMeetsSurvivalEndMargins(s)
}

// ColonistArrivalLog is the log line emitted when a colonist arrives.
const ColonistArrivalLog = colonistArrivalLog

func SnapshotDiffForTest(want, got Snapshot) string { return snapshotDiff(want, got) }

func NormalizeSnapshotForTest(s Snapshot) Snapshot { return normalizeSnapshot(s) }

func Survival45BaselineForTest() BaselineStrategy { return survival45Baseline() }
