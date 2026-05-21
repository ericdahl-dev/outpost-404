package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
	"github.com/ericdahl/outpost-404/internal/game"
)

const (
	defaultTermWidth  = 100
	defaultTermHeight = 24
)

type buildItem struct {
	id          string
	title       string
	description string
	filter      string
}

func (i buildItem) FilterValue() string { return i.filter }
func (i buildItem) Title() string      { return i.title }
func (i buildItem) Description() string  { return i.description }

func buildListItems(state game.State) []list.Item {
	items := make([]list.Item, 0, len(state.Content.Buildings))
	for _, def := range state.Content.Buildings {
		level := state.BuildingLevel(def.ID)
		cost := def.Cost * (level + 1)
		desc := fmt.Sprintf("Lv. %d/%d  Cost: %d  %s", level, def.MaxLevel, cost, def.Description)
		if level >= def.MaxLevel {
			desc = fmt.Sprintf("Lv. %d/%d  MAX  %s", level, def.MaxLevel, def.Description)
		}
		items = append(items, buildItem{
			id:          def.ID,
			title:       def.Name,
			description: desc,
			filter:      def.Name,
		})
	}
	return items
}

func newBuildListDelegate() list.DefaultDelegate {
	d := list.NewDefaultDelegate()
	d.ShowDescription = true
	d.Styles = list.NewDefaultItemStyles()
	d.Styles.SelectedTitle = d.Styles.SelectedTitle.
		Foreground(lipgloss.Color("205")).
		BorderForeground(lipgloss.Color("205"))
	d.Styles.SelectedDesc = d.Styles.SelectedDesc.
		Foreground(lipgloss.Color("241"))
	return d
}

func newBuildList(state game.State, termWidth int) list.Model {
	w := buildListWidth(termWidth)
	h := buildListHeight(len(state.Content.Buildings))
	l := list.New(buildListItems(state), newBuildListDelegate(), w, h)
	l.Title = "Build / Upgrade"
	l.SetShowTitle(true)
	l.SetShowStatusBar(false)
	l.SetShowFilter(false)
	l.SetFilteringEnabled(false)
	l.SetShowPagination(false)
	l.SetShowHelp(false)
	return l
}

func refreshBuildList(l list.Model, state game.State) list.Model {
	l.SetItems(buildListItems(state))
	return l
}

func buildListWidth(termWidth int) int {
	w := termWidth - 10
	if w > 90 {
		return 90
	}
	if w < 40 {
		return 40
	}
	return w
}

func buildListHeight(itemCount int) int {
	h := itemCount*3 + 4
	if h < 8 {
		return 8
	}
	if h > 18 {
		return 18
	}
	return h
}

func newLogViewport(termWidth, termHeight int) viewport.Model {
	vp := viewport.New(logViewportWidth(termWidth), logViewportHeight(termHeight))
	vp.Style = lipgloss.NewStyle()
	return vp
}

func syncLogViewport(vp viewport.Model, lines []string) viewport.Model {
	content := strings.Join(lines, "\n")
	if content == "" {
		content = mutedStyle.Render("Quiet shift. No new entries.")
	}
	vp.SetContent(content)
	vp.GotoBottom()
	return vp
}

func logViewportWidth(termWidth int) int {
	w := 44
	if termWidth > 120 {
		w = 52
	}
	return w
}

func logViewportHeight(termHeight int) int {
	h := termHeight - 12
	if h < 6 {
		return 6
	}
	if h > 14 {
		return 14
	}
	return h
}

func selectedBuildingID(l list.Model) (string, bool) {
	it, ok := l.SelectedItem().(buildItem)
	if !ok {
		return "", false
	}
	return it.id, true
}
