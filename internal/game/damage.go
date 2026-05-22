package game

import "fmt"

const (
	RepairBaseCost  = 25
	RepairPerLevel  = 10
)

// RepairCost is credits to clear damage on a facility at the given level.
func RepairCost(level int) int {
	if level < 1 {
		return RepairBaseCost
	}
	return RepairBaseCost + RepairPerLevel*level
}

func (s *State) damageBuilding(id string) {
	if s.BuildingLevel(id) <= 0 {
		return
	}
	b := s.Buildings[id]
	b.Damaged = true
	s.Buildings[id] = b
}

// DamageBuilding marks a built facility damaged (session-logged for replay).
func (s *State) DamageBuilding(id string) {
	s.doAction("damage", map[string]any{"building_id": id}, func(detail map[string]any) {
		s.damageWithDetail(detail, id)
	})
}

func (s *State) damageWithDetail(detail map[string]any, id string) {
	if s.GameOver {
		detail["ok"] = false
		detail["reason"] = "game_over"
		return
	}
	def, ok := s.FindBuilding(id)
	if !ok {
		s.AddLog("Unknown building.")
		detail["ok"] = false
		detail["reason"] = "unknown_building"
		return
	}
	if s.BuildingLevel(id) <= 0 {
		s.AddLog(fmt.Sprintf("%s is not built yet.", def.Name))
		detail["ok"] = false
		detail["reason"] = "not_built"
		return
	}
	s.damageBuilding(id)
	s.AddLog(fmt.Sprintf("%s took damage; daily output halved until repaired.", def.Name))
	detail["ok"] = true
	detail["building_id"] = id
}

func (s *State) firstDamagedBuildingID() string {
	for _, def := range s.Content.Buildings {
		if b, ok := s.Buildings[def.ID]; ok && b.Level > 0 && b.Damaged {
			return def.ID
		}
	}
	return ""
}

func dailyEffectsScaled(def BuildingDef, b Building) map[string]int {
	if b.Level <= 0 || len(def.DailyEffects) == 0 {
		return nil
	}
	out := make(map[string]int, len(def.DailyEffects))
	for k, v := range def.DailyEffects {
		n := v * b.Level
		if b.Damaged {
			n /= 2
		}
		if n != 0 {
			out[k] = n
		}
	}
	return out
}

func (s *State) damageRandomBuiltFacility() {
	var candidates []string
	for _, def := range s.Content.Buildings {
		if s.BuildingLevel(def.ID) > 0 {
			candidates = append(candidates, def.ID)
		}
	}
	if len(candidates) == 0 {
		return
	}
	s.ensureRNG()
	id := candidates[s.rng.Intn(len(candidates))]
	s.damageBuilding(id)
	if def, ok := s.FindBuilding(id); ok {
		s.AddLog(fmt.Sprintf("%s took damage; daily output halved until repaired.", def.Name))
	}
}
