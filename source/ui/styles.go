//
// styles.go
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

import "github.com/charmbracelet/lipgloss"

// Color palette — ANSI 256 color values used throughout the UI.
const (
	colorOrange    = lipgloss.Color("170") // titles, accents, selection
	colorPurple    = lipgloss.Color("63")  // borders, labels
	colorGray      = lipgloss.Color("241") // help text, dimmed elements
	colorLightGray = lipgloss.Color("252") // secondary values
	colorRed       = lipgloss.Color("196") // errors
	colorGreen     = lipgloss.Color("82")  // success
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorOrange).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(colorPurple).
			Padding(0, 1)

	menuItemStyle = lipgloss.NewStyle().
			PaddingLeft(4)

	selectedItemStyle = lipgloss.NewStyle().
				PaddingLeft(2).
				Foreground(colorOrange).
				Bold(true)

	helpStyle = lipgloss.NewStyle().
			Foreground(colorGray)

	outputStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(colorPurple).
			Padding(1, 2)

	errorStyle = lipgloss.NewStyle().
			Foreground(colorRed).
			Bold(true)

	successStyle = lipgloss.NewStyle().
			Foreground(colorGreen).
			Bold(true)

	progressBarWidth = 40

	progressFullStyle = lipgloss.NewStyle().
				Foreground(colorOrange)

	progressEmptyStyle = lipgloss.NewStyle().
				Foreground(colorGray)

	monitorLabelStyle = lipgloss.NewStyle().
				Foreground(colorPurple).
				Bold(true)

	monitorValueStyle = lipgloss.NewStyle().
				Foreground(colorLightGray)

	inputLabelStyle = lipgloss.NewStyle().
			Foreground(colorPurple).
			Bold(true)

	categoryStyle = lipgloss.NewStyle().
			Foreground(colorOrange).
			Bold(true)
)
