package main

import "charm.land/bubbles/v2/key"

type Keymap struct {
	Enter key.Binding
	Reset key.Binding
	Quit  key.Binding
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
	}
}

func (k Keymap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.Enter,
		k.Reset,
		k.Quit,
	}
}

func (k Keymap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Enter, k.Quit},
	}
}

func (m *model) UpdateKeymap() {
	m.keymap.Reset.SetEnabled(m.paused)
}
