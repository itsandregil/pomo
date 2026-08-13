package main

import "charm.land/bubbles/v2/key"

type Keymap struct {
	Enter key.Binding
	Quit  key.Binding
}

func DefaultKeybings() Keymap {
	return Keymap{
		Enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "resume/pause"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q", "quit"),
		),
	}
}

func (k Keymap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.Enter,
		k.Quit,
	}
}

func (k Keymap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Enter, k.Quit},
	}
}
