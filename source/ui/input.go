//
// input.go
// ~~~~~~~~~~~~~~~~~~~~~
//
// Copyright (c) 2004-2026 Metasystems Technologies Inc. (MTI)
// All rights reserved
//
// Distributed under the MTI Software License, Version 0.1.
//
// as defined by accompanying file MTI-LICENSE-0.1.info or
// at http://www.mtihq.com/license/MTI-LICENSE-0.1.info
//

package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// InputModel handles multi-field text input for parameterized commands.
type InputModel struct {
	command Command
	fields  []textinput.Model
	focus   int
	width   int
	height  int
}

// NewInputModel creates an input form for the given command.
func NewInputModel(cmd Command) InputModel {
	fields := make([]textinput.Model, len(cmd.Inputs))
	for i, inp := range cmd.Inputs {
		ti := textinput.New()
		ti.Placeholder = inp.Placeholder
		ti.CharLimit = 256
		ti.Width = 50
		if i == 0 {
			ti.Focus()
		}
		fields[i] = ti
	}
	return InputModel{
		command: cmd,
		fields:  fields,
	}
}

// Init implements tea.Model.
func (m InputModel) Init() tea.Cmd {
	return textinput.Blink
}

// inputSubmitMsg signals that the user submitted the form.
type inputSubmitMsg struct {
	command Command
	args    []string
}

// inputCancelMsg signals that the user cancelled.
type inputCancelMsg struct{}

// Update handles input form key events.
func (m InputModel) Update(msg tea.Msg) (InputModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "esc":
			return m, func() tea.Msg { return inputCancelMsg{} }
		case "tab", "down":
			return m.nextField(), nil
		case "shift+tab", "up":
			return m.prevField(), nil
		case "enter":
			// If not on last field, advance; otherwise submit
			if m.focus < len(m.fields)-1 {
				return m.nextField(), nil
			}
			return m, m.submit()
		}
	}

	// Update the focused field
	var cmd tea.Cmd
	m.fields[m.focus], cmd = m.fields[m.focus].Update(msg)
	return m, cmd
}

func (m InputModel) nextField() InputModel {
	m.fields[m.focus].Blur()
	m.focus = (m.focus + 1) % len(m.fields)
	m.fields[m.focus].Focus()
	return m
}

func (m InputModel) prevField() InputModel {
	m.fields[m.focus].Blur()
	m.focus = (m.focus - 1 + len(m.fields)) % len(m.fields)
	m.fields[m.focus].Focus()
	return m
}

func (m InputModel) submit() tea.Cmd {
	// Validate required fields
	for i, inp := range m.command.Inputs {
		if inp.Required && strings.TrimSpace(m.fields[i].Value()) == "" {
			return nil // don't submit if required fields are empty
		}
	}

	args := make([]string, len(m.fields))
	for i := range m.fields {
		args[i] = strings.TrimSpace(m.fields[i].Value())
	}
	cmd := m.command
	return func() tea.Msg {
		return inputSubmitMsg{command: cmd, args: args}
	}
}

// View renders the input form.
func (m InputModel) View() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render(m.command.Title))
	b.WriteString("\n\n")

	var form strings.Builder
	for i, inp := range m.command.Inputs {
		label := inp.Label
		if inp.Required {
			label += " *"
		}
		form.WriteString(fmt.Sprintf("%s\n", inputLabelStyle.Render(label)))
		form.WriteString(fmt.Sprintf("%s\n", m.fields[i].View()))
		if i < len(m.command.Inputs)-1 {
			form.WriteString("\n")
		}
	}
	b.WriteString(outputStyle.Render(form.String()))

	b.WriteString("\n\n")
	b.WriteString(helpStyle.Render("tab: next field • enter: submit • esc: cancel"))

	return lipgloss.Place(m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		b.String())
}
