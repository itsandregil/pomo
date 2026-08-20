package main

import "charm.land/lipgloss/v2"

// Colors
var (
	pomodoroColor   = lipgloss.Color("#FF9100")
	lightBlueColor  = lipgloss.Color("#66A3BF")
	lightGreenColor = lipgloss.Color("#66BB6A")
)

// Styles can be vars or can also be a struct
var (
	appStyle             = lipgloss.NewStyle().Margin(1, 2)
	timerBlockStyle      = lipgloss.NewStyle().Margin(1, 0)
	focusTimerStyle      = lipgloss.NewStyle().Foreground(pomodoroColor)
	shortBreakTimerStyle = lipgloss.NewStyle().Foreground(lightBlueColor)
	longTBreakimerStyle  = lipgloss.NewStyle().Foreground(lightGreenColor)
)
