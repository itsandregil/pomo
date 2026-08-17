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

type ResetTimerMsg struct{}

type model struct {
	state             state
	counter           int
	longBreakInterval int

	paused  bool
	timer   timer.Model
	spinner spinner.Model

	help     help.Model
	keymap   Keymap
	quitting bool
}

func NewModel() model {
	return model{
		state:             stateFocus,
		longBreakInterval: 4,
		paused:            true,
		timer:             timer.New(defaultFocusTime),
		spinner:           spinner.New(spinner.WithSpinner(spinner.Dot)),
		help:              help.New(),
		keymap:            DefaultKeybings(),
	}
}

func (m model) Init() tea.Cmd {
	return m.timer.Stop()
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
			m.updateKeymap()
			if !m.paused {
				return m, tea.Batch(m.timer.Start(), m.spinner.Tick)
			}
			return m, m.timer.Stop()
		case key.Matches(msg, m.keymap.Reset):
			return m, resetTimerCmd()
		case key.Matches(msg, m.keymap.Skip):
			m.state = stateFocus
			return m, resetTimerCmd()
		}
		// Toggle full help is not implemented, ignore Cmd
		m.help, _ = m.help.Update(msg)
		return m, nil
	case timer.TimeoutMsg:
		if m.state == stateFocus {
			m.counter += 1
			if m.counter >= m.longBreakInterval {
				m.counter = 0
				m.state = stateLongBreak
				return m, resetTimerCmd()
			}
			m.state = stateShortBreak
			return m, resetTimerCmd()
		}
		m.state = stateFocus
		return m, resetTimerCmd()
	case ResetTimerMsg:
		var timeout time.Duration
		switch m.state {
		case stateFocus:
			timeout = defaultFocusTime
		case stateShortBreak:
			timeout = defaultShortBreakTime
		case stateLongBreak:
			timeout = defaultLongBreakTime
		}
		m.timer = timer.New(timeout)
		m.paused = true
		m.updateKeymap()
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
		return tea.NewView("")
	}
	var message string
	switch m.state {
	case stateFocus:
		message = "Focusing..."
	case stateShortBreak, stateLongBreak:
		message = "Time for a break..."
	}
	str := "Pomo" + "\n" + m.spinner.View() + m.timer.View() + "\n" + message + "\n"
	str += fmt.Sprintf("Session %d/%d", m.counter, m.longBreakInterval)
	str += "\n" + m.help.View(m.keymap)
	return tea.NewView(str)
}

func resetTimerCmd() tea.Cmd {
	return func() tea.Msg {
		return ResetTimerMsg{}
	}
}
