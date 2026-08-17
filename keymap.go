package main

import "charm.land/bubbles/v2/key"

type Keymap struct {
	Quit  key.Binding
	Enter key.Binding
	Reset key.Binding
	Skip  key.Binding
}

func DefaultKeybings() Keymap {
	return Keymap{
		Quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q", "quit"),
		),
		Enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "start/pause"),
		),
		Reset: key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp("r", "reset"),
			key.WithDisabled(),
		),
		Skip: key.NewBinding(
			key.WithKeys("s"),
			key.WithHelp("s", "skip"),
			key.WithDisabled(),
		),
	}
}

func (k Keymap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.Enter,
		k.Reset,
		k.Skip,
		k.Quit,
	}
}

func (k Keymap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Enter, k.Quit},
	}
}

func (m *model) updateKeymap() {
	m.keymap.Reset.SetEnabled(m.paused && m.state == stateFocus)
	m.keymap.Skip.SetEnabled(m.paused && m.state != stateFocus)
}
