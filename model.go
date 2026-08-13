package main

import (
	"time"

	"charm.land/bubbles/v2/spinner"
	"charm.land/bubbles/v2/timer"
	tea "charm.land/bubbletea/v2"
)

type model struct {
	paused   bool
	timer    timer.Model
	spinner  spinner.Model
	quitting bool
}

func NewModel() model {
	t := timer.New(25 * time.Minute)
	s := spinner.New(spinner.WithSpinner(spinner.Dot))

	return model{
		timer:   t,
		spinner: s,
	}
}

func (m model) Init() tea.Cmd {
	timerCmd := m.timer.Init()
	return tea.Batch(timerCmd, m.spinner.Tick)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "q":
			m.quitting = true
			return m, tea.Quit
		case "enter":
			m.paused = !m.paused
			if !m.paused {
				return m, tea.Batch(m.timer.Start(), m.spinner.Tick)
			}
			return m, m.timer.Stop()
		}
	}
	var cmd tea.Cmd
	if !m.paused {
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)
	}
	m.timer, cmd = m.timer.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m model) View() tea.View {
	if m.quitting {
		return tea.NewView("Bye" + "\n")
	}
	str := "Pomo" + "\n" + m.spinner.View() + m.timer.View() + "\n" + "Working Session"
	return tea.NewView(str)
}
