package game

import "fmt"

func (s *State) doAction(typ string, initial map[string]any, fn func(detail map[string]any)) {
	before := s.snapshot()
	detail := initial
	if detail == nil {
		detail = map[string]any{}
	}
	fn(detail)
	s.CheckEnd()
	if !s.GameOver {
		s.syncWarnings()
	}
	s.recordAction(typ, detail, before, s.snapshot())
	if typ == "next_day" && !s.GameOver {
		_ = s.PersistAutosave()
	}
}

func (s *State) Build(id string) {
	s.doAction("build", map[string]any{"building_id": id}, func(detail map[string]any) {
		s.buildWithDetail(detail, id)
	})
}

func (s *State) buildWithDetail(detail map[string]any, id string) {
	if s.GameOver {
		detail["ok"] = false
		detail["reason"] = "game_over"
		return
	}
	def, ok := s.FindBuilding(id)
	if !ok {
		s.AddLogKind(LogSystem, "Unknown building.")
		detail["ok"] = false
		detail["reason"] = "unknown_building"
		return
	}
	level := s.BuildingLevel(id)
	if level >= def.MaxLevel {
		s.AddLogKind(LogSystem, fmt.Sprintf("%s is already at max level.", def.Name))
		detail["ok"] = false
		detail["reason"] = "max_level"
		return
	}

	cost := def.Cost * (level + 1)
	detail["cost"] = cost
	if s.Credits < cost {
		s.AddLogKind(LogSystem, fmt.Sprintf("Not enough credits for %s. Need %d.", def.Name, cost))
		detail["ok"] = false
		detail["reason"] = "insufficient_credits"
		return
	}

	s.Credits -= cost
	s.Buildings[id] = Building{DefID: id, Level: level + 1}
	s.applyEffects(def.Effects, level+1)
	newLevel := level + 1
	if newLevel == 1 {
		s.AddLogKind(LogMilestone, fmt.Sprintf("%s came online.", def.Name))
	}
	s.AddLogKind(LogGain, fmt.Sprintf("Built %s level %d. %s", def.Name, newLevel, formatEffectSummary(def.Effects)))
	s.Clamp()
	detail["ok"] = true
	detail["level"] = level + 1
}

func (s *State) Repair() {
	id := s.firstDamagedBuildingID()
	if id == "" {
		s.doAction("repair", nil, func(detail map[string]any) {
			s.AddLogKind(LogSystem, "No damaged facilities need repair.")
			detail["ok"] = false
			detail["reason"] = "nothing_damaged"
		})
		return
	}
	s.RepairBuilding(id)
}

func (s *State) RepairBuilding(id string) {
	s.doAction("repair", map[string]any{"building_id": id}, func(detail map[string]any) {
		s.repairWithDetail(detail, id)
	})
}

func (s *State) repairWithDetail(detail map[string]any, id string) {
	if s.GameOver {
		detail["ok"] = false
		detail["reason"] = "game_over"
		return
	}
	def, ok := s.FindBuilding(id)
	if !ok {
		s.AddLogKind(LogSystem, "Unknown building.")
		detail["ok"] = false
		detail["reason"] = "unknown_building"
		return
	}
	b, built := s.Buildings[id]
	if !built || b.Level <= 0 {
		s.AddLogKind(LogSystem, fmt.Sprintf("%s is not built yet.", def.Name))
		detail["ok"] = false
		detail["reason"] = "not_built"
		return
	}
	if !b.Damaged {
		s.AddLogKind(LogSystem, fmt.Sprintf("%s is not damaged.", def.Name))
		detail["ok"] = false
		detail["reason"] = "not_damaged"
		return
	}
	cost := RepairCost(b.Level)
	detail["cost"] = cost
	if s.Credits < cost {
		s.AddLogKind(LogSystem, fmt.Sprintf("Repair %s needs %d credits.", def.Name, cost))
		detail["ok"] = false
		detail["reason"] = "insufficient_credits"
		return
	}
	s.Credits -= cost
	b.Damaged = false
	s.Buildings[id] = b
	s.AddLogKind(LogGain, fmt.Sprintf("Repaired %s. Daily output restored.", def.Name))
	s.Clamp()
	detail["ok"] = true
	detail["building_id"] = id
}

func (s *State) Trade() {
	s.doAction("trade", nil, func(detail map[string]any) {
		s.tradeWithDetail(detail)
	})
}

func (s *State) tradeWithDetail(detail map[string]any) {
	if s.GameOver {
		detail["ok"] = false
		detail["reason"] = "game_over"
		return
	}
	if s.Food <= MinFoodToTrade {
		s.AddLogKind(LogSystem, fmt.Sprintf("Trade refused: food reserves too low (need more than %d).", MinFoodToTrade))
		detail["ok"] = false
		detail["reason"] = "low_food"
		return
	}
	s.Credits += TradeCreditsGain
	s.Food -= TradeFoodCost
	s.Morale -= TradeMoraleCost
	s.AddLogKind(LogTrade, fmt.Sprintf("Traded surplus rations for %d credits. Food -%d, morale -%d.", TradeCreditsGain, TradeFoodCost, TradeMoraleCost))
	s.Clamp()
	detail["ok"] = true
}

func (s *State) WorkOnBeacon() {
	s.doAction("beacon", nil, func(detail map[string]any) {
		s.beaconWithDetail(detail)
	})
}

func (s *State) beaconWithDetail(detail map[string]any) {
	if s.GameOver {
		detail["ok"] = false
		detail["reason"] = "game_over"
		return
	}
	if s.Power < 18 || s.Credits < 50 {
		s.AddLogKind(LogSystem, "Beacon work requires at least 18 power and 50 credits.")
		detail["ok"] = false
		detail["reason"] = "requirements_not_met"
		return
	}
	s.Power -= 12
	s.Credits -= 50
	s.Morale += 5
	s.BeaconParts++
	s.AddLogKind(LogMilestone, fmt.Sprintf("Signal Beacon part completed: %d/%d.", s.BeaconParts, s.MaxBeaconParts))
	s.Clamp()
	detail["ok"] = true
	detail["beacon_parts"] = s.BeaconParts
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
