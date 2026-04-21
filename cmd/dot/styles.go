package main

import (
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

// adaptive wraps a light/dark hex pair so every color works on both terminal
// backgrounds. Dark values are vivid; light values are darker shades of the
// same hue so they stay readable against a white/cream background.
func adaptive(light, dark string) lipgloss.AdaptiveColor {
	return lipgloss.AdaptiveColor{Light: light, Dark: dark}
}

var (
	// Semantic palette — used by both the huh theme and the plain lipgloss styles.
	colorPrimary = adaptive("#1d4ed8", "#4895ef") // blue   — identity / active
	colorCyan    = adaptive("#0369a1", "#48cae4") // cyan   — interactive accent
	colorSky     = adaptive("#2563eb", "#90e0ef") // sky    — descriptions, blurred
	colorSuccess = adaptive("#15803d", "#50fa7b") // green  — selected / confirmed
	colorError   = adaptive("#dc2626", "#ff5555") // red    — errors
	colorText    = adaptive("#1e293b", "#f8f8f2") // text   — normal option text
	colorMuted   = adaptive("#64748b", "#6272a4") // muted  — inactive / dim
	colorDimBg   = adaptive("#e2e8f0", "#282a36") // bg     — blurred button background
	colorAccent  = adaptive("#6d28d9", "#a78bfa") // purple — command names

	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorPrimary).
			MarginTop(1).
			MarginBottom(1)

	successStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorSuccess)

	mutedStyle = lipgloss.NewStyle().
			Foreground(colorMuted).
			Italic(true)

	commandNameStyle = lipgloss.NewStyle().
				Foreground(colorAccent).
				Bold(true)

	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorPrimary).
			Padding(0, 1)
)

const dotBanner = `
 ██████╗  ██████╗ ████████╗
 ██╔══██╗██╔═══██╗╚══██╔══╝
 ██║  ██║██║   ██║   ██║
 ██║  ██║██║   ██║   ██║
 ██████╔╝╚██████╔╝   ██║
 ╚═════╝  ╚═════╝    ╚═╝   `

// themeDot returns a custom huh theme built around the dot blue palette.
// Every color is an AdaptiveColor so the form looks sharp on both dark and
// light terminal backgrounds.
//
// Color roles at a glance:
//   - blue  — active border, titles, focused button bg
//   - cyan  — selector ›, cursor, text-input prompt
//   - sky   — descriptions, placeholder, blurred titles
//   - green — selected options, ● multiselect prefix
//   - red   — error indicators & messages
//   - text  — normal unselected option text
//   - muted — blurred / inactive text
func themeDot() *huh.Theme {
	t := huh.ThemeBase()

	var (
		blue  = colorPrimary
		cyan  = colorCyan
		sky   = colorSky
		green = colorSuccess
		red   = colorError
		text  = colorText
		muted = colorMuted
		dimBg = colorDimBg
		black = adaptive("#ffffff", "#000000") // button fg — inverts per bg
	)

	// --- Focused (active field) ---
	t.Focused.Base = t.Focused.Base.BorderForeground(blue)
	t.Focused.Card = t.Focused.Base

	t.Focused.Title = t.Focused.Title.Foreground(blue).Bold(true)
	t.Focused.NoteTitle = t.Focused.NoteTitle.Foreground(blue).Bold(true).MarginBottom(1)
	t.Focused.Directory = t.Focused.Directory.Foreground(blue)

	t.Focused.Description = t.Focused.Description.Foreground(sky)

	t.Focused.ErrorIndicator = t.Focused.ErrorIndicator.Foreground(red)
	t.Focused.ErrorMessage = t.Focused.ErrorMessage.Foreground(red)

	t.Focused.SelectSelector = t.Focused.SelectSelector.Foreground(cyan)
	t.Focused.NextIndicator = t.Focused.NextIndicator.Foreground(cyan)
	t.Focused.PrevIndicator = t.Focused.PrevIndicator.Foreground(cyan)

	t.Focused.Option = t.Focused.Option.Foreground(text)
	t.Focused.MultiSelectSelector = t.Focused.MultiSelectSelector.Foreground(cyan)
	t.Focused.SelectedOption = t.Focused.SelectedOption.Foreground(green)
	t.Focused.SelectedPrefix = lipgloss.NewStyle().Foreground(green).SetString("● ")
	t.Focused.UnselectedOption = t.Focused.UnselectedOption.Foreground(text)
	t.Focused.UnselectedPrefix = lipgloss.NewStyle().Foreground(sky).SetString("○ ")

	t.Focused.FocusedButton = t.Focused.FocusedButton.
		Foreground(black).Background(blue).Bold(true)
	t.Focused.BlurredButton = t.Focused.BlurredButton.
		Foreground(muted).Background(dimBg)
	t.Focused.Next = t.Focused.FocusedButton

	t.Focused.TextInput.Cursor = t.Focused.TextInput.Cursor.Foreground(cyan)
	t.Focused.TextInput.Placeholder = t.Focused.TextInput.Placeholder.Foreground(sky)
	t.Focused.TextInput.Prompt = t.Focused.TextInput.Prompt.Foreground(cyan)
	t.Focused.TextInput.Text = t.Focused.TextInput.Text.Foreground(text)

	// --- Blurred (inactive fields) ---
	t.Blurred = t.Focused
	t.Blurred.Base = t.Focused.Base.BorderStyle(lipgloss.HiddenBorder())
	t.Blurred.Card = t.Blurred.Base

	t.Blurred.Title = t.Blurred.Title.Foreground(sky).UnsetBold()
	t.Blurred.NoteTitle = t.Blurred.NoteTitle.Foreground(sky).UnsetBold()
	t.Blurred.Description = t.Blurred.Description.Foreground(muted)

	t.Blurred.SelectSelector = lipgloss.NewStyle().SetString("  ")
	t.Blurred.MultiSelectSelector = lipgloss.NewStyle().SetString("  ")
	t.Blurred.NextIndicator = lipgloss.NewStyle()
	t.Blurred.PrevIndicator = lipgloss.NewStyle()

	t.Blurred.TextInput.Prompt = t.Blurred.TextInput.Prompt.Foreground(sky)
	t.Blurred.TextInput.Text = t.Blurred.TextInput.Text.Foreground(muted)

	// --- Group header ---
	t.Group.Title = t.Focused.Title
	t.Group.Description = t.Focused.Description

	return t
}
