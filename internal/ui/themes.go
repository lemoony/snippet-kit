package ui

import (
	_ "embed"
	"path/filepath"
	"reflect"
	"regexp"

	"emperror.dev/errors"
	"github.com/charmbracelet/lipgloss"
	"github.com/gdamore/tcell/v2"
	"github.com/phuslu/log"
	"github.com/rivo/tview"
	"github.com/spf13/afero"
	"gopkg.in/yaml.v3"

	"github.com/lemoony/snipkit/internal/ui/style"
	"github.com/lemoony/snipkit/internal/utils/system"
	themedata "github.com/lemoony/snipkit/themes"
)

const (
	defaultThemeName       = "default"
	variablePatternMatches = 2
	filenamePatternMatches = 2
)

var (
	variablePattern = regexp.MustCompile(`^\${(?P<varName>.*)}$`)
	filenamePattern = regexp.MustCompile(`^(?P<filename>.*)\.ya?ml$`)
	ErrInvalidTheme = errors.New("invalid theme")
)

type themeWrapper struct {
	Version   string `yaml:"version"`
	Variables map[string]string
	Theme     style.ThemeValues `yaml:"theme"`
}

func applyTheme(theme style.ThemeValues) {
	styler := style.NewStyle(&theme)

	tview.Styles.PrimitiveBackgroundColor = tcell.ColorReset
	tview.Styles.BorderColor = toColor(styler.BorderColor())
	tview.Styles.TitleColor = toColor(styler.BorderTitleColor())
}

func embeddedTheme(name string) (*style.ThemeValues, bool) {
	entries, err := themedata.Files.ReadDir(".")
	if err != nil {
		panic(err)
	}

	for _, entry := range entries {
		m := filenamePattern.FindStringSubmatch(filepath.Base(entry.Name()))
		if len(m) == filenamePatternMatches {
			themeName := m[1]
			if name == themeName {
				theme := readEmbeddedTheme(entry.Name())
				return &theme, true
			}
		}
	}

	return nil, false
}

func customTheme(name string, system *system.System) (*style.ThemeValues, bool) {
	if ok, _ := afero.DirExists(system.Fs, system.ThemesDir()); !ok {
		log.Trace().Msgf("Dir does not exist: %s", system.ThemesDir())
		return nil, false
	}

	entries, err := afero.ReadDir(system.Fs, system.ThemesDir())
	if err != nil {
		panic(err)
	}

	for _, entry := range entries {
		m := filenamePattern.FindStringSubmatch(filepath.Base(entry.Name()))
		if len(m) == filenamePatternMatches {
			themeName := m[1]
			if name == themeName {
				themePath := filepath.Join(system.ThemesDir(), entry.Name())
				theme := readCustomTheme(themePath, system)
				return &theme, true
			}
		}
	}

	return nil, false
}

func readEmbeddedTheme(path string) style.ThemeValues {
	bytes, err := themedata.Files.ReadFile(path)
	if err != nil {
		panic(errors.Wrapf(err, "failed to read theme %s", path))
	}

	var wrapper themeWrapper
	err = yaml.Unmarshal(bytes, &wrapper)
	if err != nil {
		panic(err)
	}

	return wrapper.theme()
}

func readCustomTheme(path string, system *system.System) style.ThemeValues {
	bytes := system.ReadFile(path)

	var wrapper themeWrapper
	err := yaml.Unmarshal(bytes, &wrapper)
	if err != nil {
		panic(err)
	}

	return wrapper.theme()
}

func (t *themeWrapper) theme() style.ThemeValues {
	result := t.Theme
	v := reflect.Indirect(reflect.ValueOf(&result))

	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).Kind() == reflect.String {
			matches := variablePattern.FindStringSubmatch(v.Field(i).String())
			if len(matches) != variablePatternMatches {
				continue
			}
			varName := matches[1]
			if val, ok := t.Variables[varName]; !ok {
				panic(errors.Wrapf(ErrInvalidTheme, "variable %s not found", varName))
			} else {
				v.Field(i).SetString(val)
			}
		}
	}
	return result
}

func toColor(color lipgloss.TerminalColor) tcell.Color {
	if color == nil {
		return tcell.ColorReset
	}

	r, g, b, _ := color.RGBA()
	return tcell.NewRGBColor(int32(r), int32(g), int32(b))
}
