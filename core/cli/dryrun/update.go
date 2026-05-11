package dryrun

import (
	"ordin/m/core/cli/section"

	tea "charm.land/bubbletea/v2"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	if m.cursor < len(m.sections) && m.sections[m.cursor].Focus {

		updated, cmd :=
			m.sections[m.cursor].Update(msg)

		m.sections[m.cursor] =
			updated.(section.SectionModel)

		switch msg := msg.(type) {
		case tea.KeyPressMsg:
			if msg.String() == "esc" {
				m.sections[m.cursor].Focus = false
			}
		}

		return m, cmd
	} else if m.cursor == len(m.sections) && m.ambgiousSection.Focus {

		updated, cmd :=
			m.ambgiousSection.Update(msg)

		m.ambgiousSection =
			updated.(section.AmbigiousSection)

		switch msg := msg.(type) {
		case tea.KeyPressMsg:
			if msg.String() == "esc" {
				m.ambgiousSection.Focus = false
			}
		}

		return m, cmd
	} else if m.cursor == len(m.sections)+1 && m.deleteSection.Focus {

		updated, cmd :=
			m.deleteSection.Update(msg)

		m.deleteSection =
			updated.(section.DeleteSectionModel)

		switch msg := msg.(type) {
		case tea.KeyPressMsg:
			if msg.String() == "esc" {
				m.deleteSection.Focus = false
			}
		}

		return m, cmd
	}

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "tab":
			if m.cursor < (len(m.sections)-1)+2 {
				m.cursor++
			} else {
				m.cursor = 0
			}
		case "c":
			if m.cursor < len(m.sections) {
				m.sections[m.cursor].Collapse = !m.sections[m.cursor].Collapse
			} else if m.cursor == len(m.sections) {
				m.ambgiousSection.Collapse = !m.ambgiousSection.Collapse
			} else if m.cursor == len(m.sections)+1 {
				m.deleteSection.Collapse = !m.deleteSection.Collapse
			}
		case "ctrl+c", "q":
			return m, tea.Quit

		case "ctrl+s", "s":
			finalPaths := m.BuildPaths()
			SavePlan("plan.json", finalPaths)
			return m, tea.Quit

		case "enter":
			if m.cursor < len(m.sections) {
				m.sections[m.cursor].Focus = true
			} else if m.cursor == len(m.sections) {
				m.ambgiousSection.Focus = true
			}
			if m.cursor == len(m.sections)+1 {
				m.deleteSection.Focus = true
			}

		case "esc":
			if m.cursor < len(m.sections) {
				m.sections[m.cursor].Focus = false
			} else if m.cursor == len(m.sections) {
				m.ambgiousSection.Focus = false
			}
			if m.cursor == len(m.sections)+1 {
				m.deleteSection.Focus = false
			}

		}

	}
	return m, nil
}
