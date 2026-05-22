package game

import "math/rand"

func (s *State) rngIntn(n int) int {
	s.ensureRNG()
	s.rngDrawMods = append(s.rngDrawMods, n)
	return s.rng.Intn(n)
}

func replayRNG(seed int64, drawMods []int) *rand.Rand {
	r := rand.New(rand.NewSource(seed))
	for _, n := range drawMods {
		if n > 0 {
			r.Intn(n)
		}
	}
	return r
}
