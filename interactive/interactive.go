package interactive

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
)

type errMsg error

type Model struct {
	textarea textarea.Model
	err      error
}

func InitialModel() Model {
	ti := textarea.New()
	ti.Placeholder = "<paste tree command output>"
	ti.Focus()

	return Model{
		textarea: ti,
		err:      nil,
	}
}

func (m Model) Init() tea.Cmd {
	return textarea.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			fmt.Println(m.textarea.Value())
			return m, tea.Quit
		default:
			if !m.textarea.Focused() {
				cmd = m.textarea.Focus()
				cmds = append(cmds, cmd)
			}
		}

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	m.textarea, cmd = m.textarea.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	return fmt.Sprintf(
		"Insert the tree:\n\n%s\n\n%s",
		m.textarea.View(),
		"Save ([ctrl] + [c])",
	) + "\n\n"
}
