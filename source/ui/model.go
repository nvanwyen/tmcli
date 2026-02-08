//
// model.go
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

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type viewState int

const (
	categoryView viewState = iota
	commandView
	outputView
	monitorView
	inputView
	helpCategoryView
	helpCommandView
	helpDetailView
)

type commandResultMsg struct {
	output string
	err    error
}

// Model is the top-level Bubbletea model.
type Model struct {
	view       viewState
	categories []Category
	catCursor  int // cursor within category menu
	cmdCursor  int // cursor within command submenu
	output       string
	scrollOffset int
	err          error
	width      int
	height     int
	monitor      MonitorModel
	input        InputModel
	helpCursor    int    // cursor within help category picker
	helpCmdCursor int    // cursor within help command list
	helpOutput    string // rendered help text for detail view
}

// NewModel returns the initial model.
func NewModel() Model {
	return Model{
		view:       categoryView,
		categories: Categories(),
	}
}

// Init implements tea.Model.
func (m Model) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch m.view {
		case categoryView:
			return m.updateCategory(msg)
		case commandView:
			return m.updateCommand(msg)
		case outputView:
			return m.updateOutput(msg)
		case monitorView:
			return m.updateMonitor(msg)
		case inputView:
			return m.updateInput(msg)
		case helpCategoryView:
			return m.updateHelpCategory(msg)
		case helpCommandView:
			return m.updateHelpCommand(msg)
		case helpDetailView:
			return m.updateHelpDetail(msg)
		}

	case statusUpdateMsg, statusTickMsg:
		if m.view == monitorView {
			return m.updateMonitor(msg)
		}

	case commandResultMsg:
		m.output = msg.output
		m.err = msg.err
		m.scrollOffset = 0
		m.view = outputView
		return m, nil

	case inputSubmitMsg:
		m.view = outputView
		return m, m.executeWithArgs(msg.command, msg.args)

	case inputCancelMsg:
		m.view = commandView
		return m, nil
	}

	return m, nil
}

// --- Category menu ---

func (m Model) updateCategory(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	helpIdx := len(m.categories)
	quitIdx := helpIdx + 1
	count := quitIdx + 1
	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit
	case "up", "k":
		if m.catCursor > 0 {
			m.catCursor--
		} else {
			m.catCursor = count - 1
		}
	case "down", "j":
		if m.catCursor < count-1 {
			m.catCursor++
		} else {
			m.catCursor = 0
		}
	case "enter":
		return m.selectCategoryItem()
	case "h":
		m.view = helpCategoryView
		m.helpCursor = 0
		return m, nil
	case "q":
		return m, tea.Quit
	default:
		// Check category hotkeys
		for i, cat := range m.categories {
			if msg.String() == cat.Hotkey {
				m.catCursor = i
				m.view = commandView
				m.cmdCursor = 0
				return m, nil
			}
		}
	}
	return m, nil
}

func (m Model) selectCategoryItem() (tea.Model, tea.Cmd) {
	helpIdx := len(m.categories)
	quitIdx := helpIdx + 1
	switch m.catCursor {
	case quitIdx:
		return m, tea.Quit
	case helpIdx:
		m.view = helpCategoryView
		m.helpCursor = 0
		return m, nil
	default:
		m.view = commandView
		m.cmdCursor = 0
		return m, nil
	}
}

// --- Command submenu ---

func (m Model) updateCommand(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	cmds := m.categories[m.catCursor].Commands
	backIdx := len(cmds)
	quitIdx := backIdx + 1
	count := quitIdx + 1
	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit
	case "esc", "backspace":
		m.view = categoryView
		return m, nil
	case "up", "k":
		if m.cmdCursor > 0 {
			m.cmdCursor--
		} else {
			m.cmdCursor = count - 1
		}
	case "down", "j":
		if m.cmdCursor < count-1 {
			m.cmdCursor++
		} else {
			m.cmdCursor = 0
		}
	case "enter":
		switch m.cmdCursor {
		case quitIdx:
			return m, tea.Quit
		case backIdx:
			m.view = categoryView
			return m, nil
		default:
			return m.selectCommand(cmds[m.cmdCursor])
		}
	default:
		// Check command hotkeys
		for _, cmd := range cmds {
			if msg.String() == cmd.Hotkey {
				return m.selectCommand(cmd)
			}
		}
		if msg.String() == "b" {
			m.view = categoryView
			return m, nil
		}
		if msg.String() == "q" {
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m Model) selectCommand(cmd Command) (tea.Model, tea.Cmd) {
	if cmd.IsMonitor {
		m.monitor = NewMonitorModel(true)
		m.view = monitorView
		return m, m.monitor.Init()
	}
	if len(cmd.Inputs) > 0 {
		m.input = NewInputModel(cmd)
		m.input.width = m.width
		m.input.height = m.height
		m.view = inputView
		return m, m.input.Init()
	}
	m.view = outputView
	return m, m.executeWithArgs(cmd, nil)
}

func (m Model) executeWithArgs(cmd Command, args []string) tea.Cmd {
	return func() tea.Msg {
		output, err := cmd.Execute(args)
		return commandResultMsg{output: output, err: err}
	}
}

// --- Output view ---

func (m Model) updateOutput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit
	case "q":
		return m, tea.Quit
	case "esc", "backspace", "b":
		m.view = commandView
		m.output = ""
		m.err = nil
		m.scrollOffset = 0
	case "up", "k":
		if m.scrollOffset > 0 {
			m.scrollOffset--
		}
	case "down", "j":
		lines := strings.Split(m.output, "\n")
		maxOff := len(lines) - m.outputPageSize()
		if maxOff < 0 {
			maxOff = 0
		}
		if m.scrollOffset < maxOff {
			m.scrollOffset++
		}
	case "pgup":
		m.scrollOffset -= m.outputPageSize()
		if m.scrollOffset < 0 {
			m.scrollOffset = 0
		}
	case "pgdown", " ":
		lines := strings.Split(m.output, "\n")
		maxOff := len(lines) - m.outputPageSize()
		if maxOff < 0 {
			maxOff = 0
		}
		m.scrollOffset += m.outputPageSize()
		if m.scrollOffset > maxOff {
			m.scrollOffset = maxOff
		}
	}
	return m, nil
}

func (m Model) outputPageSize() int {
	ps := m.height - 12
	if ps < 5 {
		ps = 5
	}
	return ps
}

// --- Monitor view ---

func (m Model) updateMonitor(msg tea.Msg) (tea.Model, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.String() {
		case "esc", "backspace", "b":
			m.view = commandView
			return m, nil
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}

	updated, cmd := m.monitor.Update(msg)
	m.monitor = updated.(MonitorModel)
	return m, cmd
}

// --- Input view ---

func (m Model) updateInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	updated, cmd := m.input.Update(msg)
	m.input = updated
	return m, cmd
}

// --- Help views ---

func (m Model) updateHelpCategory(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	cats := Categories()
	backIdx := len(cats)
	quitIdx := backIdx + 1
	count := quitIdx + 1
	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit
	case "esc", "backspace":
		m.view = categoryView
		return m, nil
	case "up", "k":
		if m.helpCursor > 0 {
			m.helpCursor--
		} else {
			m.helpCursor = count - 1
		}
	case "down", "j":
		if m.helpCursor < count-1 {
			m.helpCursor++
		} else {
			m.helpCursor = 0
		}
	case "enter":
		switch m.helpCursor {
		case quitIdx:
			return m, tea.Quit
		case backIdx:
			m.view = categoryView
			return m, nil
		default:
			m.view = helpCommandView
			m.helpCmdCursor = 0
			return m, nil
		}
	case "b":
		m.view = categoryView
		return m, nil
	case "q":
		return m, tea.Quit
	default:
		for i, cat := range cats {
			if msg.String() == cat.Hotkey {
				m.helpCursor = i
				m.view = helpCommandView
				m.helpCmdCursor = 0
				return m, nil
			}
		}
	}
	return m, nil
}

func (m Model) updateHelpCommand(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	cats := Categories()
	cmds := cats[m.helpCursor].Commands
	backIdx := len(cmds)
	quitIdx := backIdx + 1
	count := quitIdx + 1
	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit
	case "esc", "backspace":
		m.view = helpCategoryView
		return m, nil
	case "up", "k":
		if m.helpCmdCursor > 0 {
			m.helpCmdCursor--
		} else {
			m.helpCmdCursor = count - 1
		}
	case "down", "j":
		if m.helpCmdCursor < count-1 {
			m.helpCmdCursor++
		} else {
			m.helpCmdCursor = 0
		}
	case "enter":
		switch m.helpCmdCursor {
		case quitIdx:
			return m, tea.Quit
		case backIdx:
			m.view = helpCategoryView
			return m, nil
		default:
			m.helpOutput = BuildCommandHelp(cmds[m.helpCmdCursor])
			m.view = helpDetailView
			return m, nil
		}
	case "b":
		m.view = helpCategoryView
		return m, nil
	case "q":
		return m, tea.Quit
	default:
		for i, cmd := range cmds {
			if msg.String() == cmd.Hotkey {
				m.helpOutput = BuildCommandHelp(cmds[i])
				m.view = helpDetailView
				return m, nil
			}
		}
	}
	return m, nil
}

func (m Model) updateHelpDetail(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit
	case "esc", "backspace", "b":
		m.view = helpCommandView
		return m, nil
	case "q":
		return m, tea.Quit
	}
	return m, nil
}

// --- Views ---

// View implements tea.Model.
func (m Model) View() string {
	switch m.view {
	case categoryView:
		return m.renderCategory()
	case commandView:
		return m.renderCommand()
	case outputView:
		return m.renderOutput()
	case monitorView:
		return m.renderMonitor()
	case inputView:
		return m.input.View()
	case helpCategoryView:
		return m.renderHelpCategory()
	case helpCommandView:
		return m.renderHelpCommand()
	case helpDetailView:
		return m.renderHelpDetail()
	}
	return ""
}

func (m Model) renderCategory() string {
	var b strings.Builder

	title := titleStyle.Render("Time Machine CLI")
	b.WriteString(title)
	b.WriteString("\n\n")

	helpIdx := len(m.categories)
	quitIdx := helpIdx + 1

	var menu strings.Builder
	for i, cat := range m.categories {
		if i == m.catCursor {
			fmt.Fprintf(&menu, "> [%s] %s\n", cat.Hotkey, cat.Title)
		} else {
			fmt.Fprintf(&menu, "  [%s] %s\n", cat.Hotkey, cat.Title)
		}
	}
	menu.WriteString("\n")
	if m.catCursor == helpIdx {
		fmt.Fprintf(&menu, "> [h] Help\n")
	} else {
		fmt.Fprintf(&menu, "  [h] Help\n")
	}
	if m.catCursor == quitIdx {
		fmt.Fprintf(&menu, "> [q] Quit\n")
	} else {
		fmt.Fprintf(&menu, "  [q] Quit\n")
	}
	b.WriteString(outputStyle.Render(menu.String()))

	b.WriteString("\n\n")
	b.WriteString(helpStyle.Render("↑/↓: navigate • enter/hotkey: select"))

	return lipgloss.Place(m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		b.String())
}

func (m Model) renderCommand() string {
	var b strings.Builder

	cat := m.categories[m.catCursor]
	title := titleStyle.Render(cat.Title)
	b.WriteString(title)
	b.WriteString("\n\n")

	// Check if any command in this category requires root
	hasRoot := false
	for _, cmd := range cat.Commands {
		if cmd.RequiresRoot {
			hasRoot = true
			break
		}
	}

	// Calculate max item width for right-aligned asterisks
	maxW := len("  [b] Back")
	for _, cmd := range cat.Commands {
		w := len(fmt.Sprintf("  [%s] %s", cmd.Hotkey, cmd.Title))
		if w > maxW {
			maxW = w
		}
	}

	var menu strings.Builder
	for i, cmd := range cat.Commands {
		var line string
		if i == m.cmdCursor {
			line = fmt.Sprintf("> [%s] %s", cmd.Hotkey, cmd.Title)
		} else {
			line = fmt.Sprintf("  [%s] %s", cmd.Hotkey, cmd.Title)
		}
		if hasRoot {
			if cmd.RequiresRoot {
				fmt.Fprintf(&menu, "%-*s *\n", maxW, line)
			} else {
				fmt.Fprintf(&menu, "%-*s  \n", maxW, line)
			}
		} else {
			fmt.Fprintf(&menu, "%s\n", line)
		}
	}
	menu.WriteString("\n")
	backLine := "  [b] Back"
	if m.cmdCursor == len(cat.Commands) {
		backLine = "> [b] Back"
	}
	quitLine := "  [q] Quit"
	if m.cmdCursor == len(cat.Commands)+1 {
		quitLine = "> [q] Quit"
	}
	if hasRoot {
		fmt.Fprintf(&menu, "%-*s  \n", maxW, backLine)
		fmt.Fprintf(&menu, "%-*s  \n", maxW, quitLine)
	} else {
		fmt.Fprintf(&menu, "%s\n", backLine)
		fmt.Fprintf(&menu, "%s\n", quitLine)
	}
	b.WriteString(outputStyle.Render(menu.String()))

	if hasRoot {
		b.WriteString("\n")
		b.WriteString(helpStyle.Render("* requires root"))
	}

	b.WriteString("\n\n")
	b.WriteString(helpStyle.Render("↑/↓: navigate • enter/hotkey: select • esc: back"))

	return lipgloss.Place(m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		b.String())
}

func (m Model) renderOutput() string {
	var b strings.Builder

	title := titleStyle.Render("Time Machine CLI")
	b.WriteString(title)
	b.WriteString("\n\n")

	if m.err != nil {
		b.WriteString(errorStyle.Render(fmt.Sprintf("Error: %v", m.err)))
		b.WriteString("\n\n")
		b.WriteString(helpStyle.Render("b/esc: back • q: quit"))
	} else {
		lines := strings.Split(m.output, "\n")
		pageSize := m.outputPageSize()

		if len(lines) <= pageSize {
			b.WriteString(outputStyle.Render(m.output))
			b.WriteString("\n\n")
			b.WriteString(helpStyle.Render("b/esc: back • q: quit"))
		} else {
			end := m.scrollOffset + pageSize
			if end > len(lines) {
				end = len(lines)
			}
			page := strings.Join(lines[m.scrollOffset:end], "\n")
			b.WriteString(outputStyle.Render(page))
			b.WriteString("\n\n")
			b.WriteString(helpStyle.Render(
				fmt.Sprintf("↑/↓: scroll • pgup/pgdn: page • lines %d–%d of %d • b/esc: back • q: quit",
					m.scrollOffset+1, end, len(lines))))
		}
	}

	return lipgloss.Place(m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		b.String())
}

func (m Model) renderMonitor() string {
	m.monitor.width = m.width
	m.monitor.height = m.height
	return m.monitor.View()
}

func (m Model) renderHelpCategory() string {
	var b strings.Builder

	title := titleStyle.Render("Help")
	b.WriteString(title)
	b.WriteString("\n\n")

	cats := Categories()
	var menu strings.Builder
	for i, cat := range cats {
		if i == m.helpCursor {
			fmt.Fprintf(&menu, "> [%s] %s\n", cat.Hotkey, cat.Title)
		} else {
			fmt.Fprintf(&menu, "  [%s] %s\n", cat.Hotkey, cat.Title)
		}
	}
	menu.WriteString("\n")
	backIdx := len(cats)
	quitIdx := backIdx + 1
	if m.helpCursor == backIdx {
		fmt.Fprintf(&menu, "> [b] Back\n")
	} else {
		fmt.Fprintf(&menu, "  [b] Back\n")
	}
	if m.helpCursor == quitIdx {
		fmt.Fprintf(&menu, "> [q] Quit\n")
	} else {
		fmt.Fprintf(&menu, "  [q] Quit\n")
	}
	b.WriteString(outputStyle.Render(menu.String()))

	b.WriteString("\n\n")
	b.WriteString(helpStyle.Render("↑/↓: navigate • enter/hotkey: select • esc: back"))

	return lipgloss.Place(m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		b.String())
}

func (m Model) renderHelpCommand() string {
	var b strings.Builder

	cats := Categories()
	cat := cats[m.helpCursor]
	title := titleStyle.Render("Help — " + cat.Title)
	b.WriteString(title)
	b.WriteString("\n\n")

	// Check if any command in this category requires root
	hasRoot := false
	for _, cmd := range cat.Commands {
		if cmd.RequiresRoot {
			hasRoot = true
			break
		}
	}

	// Calculate max item width for right-aligned asterisks
	maxW := len("  [b] Back")
	for _, cmd := range cat.Commands {
		w := len(fmt.Sprintf("  [%s] %s", cmd.Hotkey, cmd.Title))
		if w > maxW {
			maxW = w
		}
	}

	var menu strings.Builder
	for i, cmd := range cat.Commands {
		var line string
		if i == m.helpCmdCursor {
			line = fmt.Sprintf("> [%s] %s", cmd.Hotkey, cmd.Title)
		} else {
			line = fmt.Sprintf("  [%s] %s", cmd.Hotkey, cmd.Title)
		}
		if hasRoot {
			if cmd.RequiresRoot {
				fmt.Fprintf(&menu, "%-*s *\n", maxW, line)
			} else {
				fmt.Fprintf(&menu, "%-*s  \n", maxW, line)
			}
		} else {
			fmt.Fprintf(&menu, "%s\n", line)
		}
	}
	menu.WriteString("\n")
	backLine := "  [b] Back"
	if m.helpCmdCursor == len(cat.Commands) {
		backLine = "> [b] Back"
	}
	quitLine := "  [q] Quit"
	if m.helpCmdCursor == len(cat.Commands)+1 {
		quitLine = "> [q] Quit"
	}
	if hasRoot {
		fmt.Fprintf(&menu, "%-*s  \n", maxW, backLine)
		fmt.Fprintf(&menu, "%-*s  \n", maxW, quitLine)
	} else {
		fmt.Fprintf(&menu, "%s\n", backLine)
		fmt.Fprintf(&menu, "%s\n", quitLine)
	}
	b.WriteString(outputStyle.Render(menu.String()))

	if hasRoot {
		b.WriteString("\n")
		b.WriteString(helpStyle.Render("* requires root"))
	}

	b.WriteString("\n\n")
	b.WriteString(helpStyle.Render("↑/↓: navigate • enter/hotkey: select • esc: back"))

	return lipgloss.Place(m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		b.String())
}

func (m Model) renderHelpDetail() string {
	var b strings.Builder

	title := titleStyle.Render("Help")
	b.WriteString(title)
	b.WriteString("\n\n")
	b.WriteString(outputStyle.Render(m.helpOutput))
	b.WriteString("\n\n")
	b.WriteString(helpStyle.Render("b/esc: back • q: quit"))

	return lipgloss.Place(m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		b.String())
}
