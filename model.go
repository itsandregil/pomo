package main

import (
	"fmt"
	"time"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/spinner"
	"charm.land/bubbles/v2/timer"
	tea "charm.land/bubbletea/v2"
)

type state int

const (
	stateFocus state = iota
	stateShortBreak
	stateLongBreak
)

const (
	defaultFocusTime      = 25 * time.Minute
	defaultShortBreakTime = 5 * time.Minute
	defaultLongBreakTime  = 15 * time.Minute
)

type model struct {
	state              state
	pomoCounter        int
	cyclesForLongBreak int

	paused  bool
	timer   timer.Model
	spinner spinner.Model

	help     help.Model
	keymap   Keymap
	quitting bool
}

func NewModel() model {
	return model{
		state:              stateFocus,
		cyclesForLongBreak: 4,
		timer:              timer.New(defaultFocusTime),
		spinner:            spinner.New(spinner.WithSpinner(spinner.Dot)),
		help:               help.New(),
		keymap:             DefaultKeybings(),
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
		switch {
		case key.Matches(msg, m.keymap.Quit):
			m.quitting = true
			return m, tea.Quit
		case key.Matches(msg, m.keymap.Enter):
			m.paused = !m.paused
			m.UpdateKeymap()
			if !m.paused {
				return m, tea.Batch(m.timer.Start(), m.spinner.Tick)
			}
			return m, m.timer.Stop()
		case key.Matches(msg, m.keymap.Reset):
			m.resetTimer()
			return m, nil
		}
		m.help, _ = m.help.Update(msg)
		return m, nil
	case timer.TimeoutMsg:
		switch m.state {
		case stateFocus:
			m.pomoCounter += 1
			if m.pomoCounter >= m.cyclesForLongBreak {
				m.state = stateLongBreak
			} else {
				m.state = stateShortBreak
			}
		case stateShortBreak:
			m.state = stateFocus
		case stateLongBreak:
			m.state = stateFocus
			m.pomoCounter = 0
		}
		m.resetTimer()
		return m, nil
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
	str := "Pomo" + "\n" + m.spinner.View() + m.timer.View() + "\n" + m.sessionView() + "\n"
	str += fmt.Sprintf("Session %d/%d", m.pomoCounter, m.cyclesForLongBreak)
	str += "\n" + m.help.View(m.keymap)
	return tea.NewView(str)
}

func (m *model) resetTimer() {
	m.timer = timer.New(m.nextInterval())
	m.paused = true
}

func (m model) nextInterval() time.Duration {
	if m.state == stateShortBreak {
		return defaultShortBreakTime
	}
	if m.state == stateLongBreak {
		return defaultLongBreakTime
	}
	return defaultFocusTime
}

func (m model) sessionView() string {
	if m.state == stateShortBreak {
		return "Time for a break..."
	}
	if m.state == stateLongBreak {
		return "Take a break! You deserve it."
	}
	return "Focusing..."
}
