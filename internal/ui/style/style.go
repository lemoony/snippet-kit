package style

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

type Style struct {
	width  int
	height int

	colors   colors
	showHelp bool
	minimize bool

	needsToResize bool

	Reload func()
}

var NoopStyle = &Style{}

func NewStyle(t *ThemeValues, showKeyMap bool) Style {
	return Style{
		colors:   newColors(t),
		showHelp: showKeyMap,
	}
}

func (s *Style) SetSize(width, height int) {
	s.width = width
	s.height = height
}

func (s *Style) NeedsResize() bool {
	return s.needsToResize
}

func (s *Style) Profile() termenv.Profile {
	return colorProfile
}

func (s *Style) TitleStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Background(s.colors.titleColor.Value()).
		Foreground(s.colors.titleContrastColor.Value()).
		Bold(true).
		Italic(true).
		Padding(0, 1).
		MarginBottom(1)
}

func (s *Style) Title(text string) string {
	return s.TitleStyle().Render(text)
}

func (s *Style) PromptLabelStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(s.TitleColor().Value()).
		Border(lipgloss.ThickBorder(), false, false, false, true).
		BorderForeground(s.BorderColor().Value()).
		PaddingLeft(1)
}

func (s *Style) PromptLabel(text string) string {
	return s.PromptLabelStyle().Render(text)
}

func (s *Style) PromptDescriptionStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Italic(true).
		Foreground(s.PlaceholderColor().Value()).
		Border(lipgloss.ThickBorder(), false, false, false, true).
		BorderForeground(s.BorderColor().Value()).
		PaddingLeft(1)
}

func (s *Style) PromptDescription(text string) string {
	return s.PromptDescriptionStyle().Render(text)
}

func (s *Style) InputIndentSyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Border(lipgloss.ThickBorder(), false, false, false, true).
		BorderForeground(s.BorderColor().Value()).
		PaddingLeft(1)
}

func (s *Style) InputIndent(text string) string {
	return s.InputIndentSyle().Render(text)
}

func (s *Style) InputHelpStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(s.PlaceholderColor().Value())
}

func (s *Style) InputHelp(text string) string {
	return s.InputHelpStyle().Render(text)
}

func (s *Style) FormFieldWrapper(field string) string {
	if s.minimize {
		return lipgloss.NewStyle().Margin(0, 0, 0, 0).Render(field)
	}
	return lipgloss.NewStyle().Margin(0, 0, 1, 0).Render(field)
}

func (s *Style) MainView(view string, help string, resize bool) string {
	defer func() {
		if resize {
			s.needsToResize = false
		}
	}()

	var sections []string
	sections = append(sections, view)

	marginsDefault := []int{1, 2, 1, 4}
	marginsMinimal := []int{0, 2, 0, 4}

	viewHeight := lipgloss.Height(view)
	helpHeight := 0
	if s.showHelp {
		helpHeight = lipgloss.Height(help)
	}

	var margins []int
	if viewHeight+helpHeight+marginsDefault[0]+marginsDefault[2] <= s.height {
		margins = marginsDefault
		if s.minimize {
			s.minimize = false
			if !resize {
				s.needsToResize = true
			}
		}
	} else {
		margins = marginsMinimal
		if !s.minimize {
			s.minimize = true
			if !resize {
				s.needsToResize = true
			}
		}
	}

	// Fill empty space with newlines
	extraLines := ""
	if requiredHeight := viewHeight + helpHeight + margins[0] + margins[2]; requiredHeight < s.height {
		extraLines = strings.Repeat("\n", max(0, s.height-requiredHeight-1))
	}

	if extraLines != "" {
		sections = append(sections, extraLines)
	}

	if s.showHelp {
		sections = append(sections, help)
	}

	return lipgloss.NewStyle().Margin(margins...).Render(lipgloss.JoinVertical(lipgloss.Left, sections...))
}

func (s Style) ColorProfile() termenv.Profile {
	return colorProfile
}

func (s Style) PreviewColorSchemeName() string {
	return s.colors.previewColorSchemeName
}

func (s Style) BorderColor() Color {
	return s.colors.borderColor
}

func (s Style) BorderTitleColor() Color {
	return s.colors.borderColor
}

func (s Style) TitleColor() Color {
	return s.colors.titleColor
}

func (s Style) TitleContrastColor() Color {
	return s.colors.titleContrastColor
}

func (s Style) TextColor() Color {
	return s.colors.textColor
}

func (s Style) PlaceholderColor() Color {
	return s.colors.subduedColor
}

func (s Style) SubduedColor() Color {
	return s.colors.subduedColor
}

func (s Style) VerySubduedColor() Color {
	return s.colors.verySubduedColor
}

func (s Style) ActiveColor() Color {
	return s.colors.activeColor
}

func (s Style) ActiveContrastColor() Color {
	return s.colors.activeContrastColor
}

func (s Style) InfoColor() Color {
	return s.colors.infoColor
}

func (s Style) HighlightColor() Color {
	return s.colors.highlightColor
}

func (s Style) HighlightContrastColor() Color {
	return s.colors.highlightContrastColor
}

func (s Style) SnippetColor() Color {
	return s.colors.snippetColor
}

func (s Style) SnippetContrastColor() Color {
	return s.colors.snippetContrastColor
}

func (s Style) ButtonTextColor(selected bool) Color {
	if selected {
		return s.colors.activeContrastColor
	}
	return s.colors.subduedContrastColor
}

func (s Style) ButtonColor(selected bool) Color {
	if selected {
		return s.colors.activeColor
	}
	return s.colors.subduedColor
}

func (s Style) SuccessColor() Color {
	return s.colors.successColor
}

func (s Style) ErrorColor() Color {
	return s.colors.errorColor
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
