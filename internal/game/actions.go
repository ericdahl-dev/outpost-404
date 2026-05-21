package game

import "fmt"

func (s *State) doAction(typ string, initial map[string]any, fn func(detail map[string]any)) {
	before := s.snapshot()
	detail := initial
	if detail == nil {
		detail = map[string]any{}
	}
	fn(detail)
	s.recordAction(typ, detail, before, s.snapshot())
	s.CheckEnd()
}

func (s *State) Build(id string) {
	s.doAction("build", map[string]any{"building_id": id}, func(detail map[string]any) {
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
		level := s.BuildingLevel(id)
		if level >= def.MaxLevel {
			s.AddLog(fmt.Sprintf("%s is already at max level.", def.Name))
			detail["ok"] = false
			detail["reason"] = "max_level"
			return
		}

		cost := def.Cost * (level + 1)
		detail["cost"] = cost
		if s.Credits < cost {
			s.AddLog(fmt.Sprintf("Not enough credits for %s. Need %d.", def.Name, cost))
			detail["ok"] = false
			detail["reason"] = "insufficient_credits"
			return
		}

		s.Credits -= cost
		s.Buildings[id] = Building{DefID: id, Level: level + 1}
		s.applyEffects(def.Effects, level+1)
		s.AddLog(fmt.Sprintf("Built %s level %d.", def.Name, level+1))
		s.Clamp()
		detail["ok"] = true
		detail["level"] = level + 1
	})
}

func (s *State) Repair() {
	s.doAction("repair", nil, func(detail map[string]any) {
		if s.GameOver {
			detail["ok"] = false
			detail["reason"] = "game_over"
			return
		}
		if s.Credits < 35 {
			s.AddLog("Repairs require 35 credits.")
			detail["ok"] = false
			detail["reason"] = "insufficient_credits"
			return
		}
		s.Credits -= 35
		s.Power += 12
		s.Morale += 4
		s.AddLog("Workshop crew patched failing systems. Power +12, morale +4.")
		s.Clamp()
		detail["ok"] = true
	})
}

func (s *State) Trade() {
	s.doAction("trade", nil, func(detail map[string]any) {
		if s.GameOver {
			detail["ok"] = false
			detail["reason"] = "game_over"
			return
		}
		s.Credits += 45
		s.Food -= 8
		s.Morale -= 3
		s.AddLog("Traded surplus rations for 45 credits. Food -8, morale -3.")
		s.Clamp()
		detail["ok"] = true
	})
}

func (s *State) WorkOnBeacon() {
	s.doAction("beacon", nil, func(detail map[string]any) {
		if s.GameOver {
			detail["ok"] = false
			detail["reason"] = "game_over"
			return
		}
		if s.Power < 18 || s.Credits < 50 {
			s.AddLog("Beacon work requires at least 18 power and 50 credits.")
			detail["ok"] = false
			detail["reason"] = "requirements_not_met"
			return
		}
		s.Power -= 12
		s.Credits -= 50
		s.Morale += 5
		s.BeaconParts++
		s.AddLog(fmt.Sprintf("Signal Beacon part completed: %d/%d.", s.BeaconParts, s.MaxBeaconParts))
		s.Clamp()
		detail["ok"] = true
		detail["beacon_parts"] = s.BeaconParts
	})
}

func (s *State) applyEffects(effects map[string]int, multiplier int) {
	for key, amount := range effects {
		amount *= multiplier
		switch Resource(key) {
		case ResourcePower:
			s.Power += amount
		case ResourceFood:
			s.Food += amount
		case ResourceMorale:
			s.Morale += amount
		case ResourceCredits:
			s.Credits += amount
		default:
			switch key {
			case "populationCap":
				s.PopulationCap += amount
			case "population":
				s.Population += amount
			}
		}
	}
}
