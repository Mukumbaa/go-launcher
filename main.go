package main

import (
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	// "github.com/charmbracelet/lipgloss"
)


var TERMINAL string = "alacritty"
var BROWSER string = "firefox"
var THEME string = ""
var theme Theme

type model struct {
	textInput    textinput.Model
	apps         []AppEntry
	filteredApps []AppEntry
	cursor       int
	offset       int
	selectedApp  *AppEntry
	
    width  int
    height int
}

func initialModel() model {
	loadConfig()
	
	palette, err := loadThemeConfig(THEME)
	if err != nil {
		fmt.Printf("error reading theme config: %v", err)
	}

	theme = BuildTheme(palette)
	
	ti := textinput.New()
	ti.Placeholder = "Search..."
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 50

	ti.TextStyle = theme.Input.Text
	ti.PlaceholderStyle = theme.Input.Placeholder
	ti.PromptStyle = theme.Input.Prompt
	
	args := os.Args[1:]

	var apps []AppEntry
	if len(args) > 0 {
	    apps = loadFromFile(args[0])
	} else {
	    apps = loadApplications()
	}


	time.Sleep(10*time.Millisecond)

	return model{
		textInput:    ti,
		apps:         apps,
		filteredApps: apps,
		cursor:       0,
		offset:       0,
		selectedApp:  nil,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	oldValue := m.textInput.Value()

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit

		case "enter":
			if len(m.filteredApps) > 0 && m.cursor < len(m.filteredApps) {
				selected := m.filteredApps[m.cursor]
				m.selectedApp = &selected
				return m, tea.Quit
			}

		case "up", "shift+tab":
			if m.cursor > 0 {
				m.cursor--
				if m.cursor < m.offset {
					m.offset--
				}
			}

		case "down", "tab":
			if m.cursor < len(m.filteredApps)-1 {
				m.cursor++
				if m.cursor >= m.offset+10{
					m.offset++
				}
			}
		}

	m.textInput, cmd = m.textInput.Update(msg)

	rawInput := m.textInput.Value() 
	if oldValue != rawInput {
		m.cursor = 0
		m.offset = 0
	}

	
	query := strings.TrimSpace(rawInput)
	
	if strings.HasPrefix(rawInput, "?") {
		searchText := strings.TrimSpace(strings.TrimPrefix(rawInput, "?"))
		
		var displayName string
		var execCmd string
		
		if searchText == "" {
			displayName = "Google search:"
			execCmd = ""
			m.filteredApps = []AppEntry{} 
		} else {
			encodedQuery := url.QueryEscape(searchText)
			
			googleUrl := fmt.Sprintf("https://www.google.com/search?q=%s", encodedQuery)
			displayName = fmt.Sprintf("Google search: %s", searchText)
			
			execCmd = fmt.Sprintf("%s '%s'", BROWSER,googleUrl)
			
		}
		webSearchEntry := AppEntry{
			Name:     displayName,
			Exec:     execCmd,
			Terminal: false,
			File:     "", // Non esiste un file .desktop reale
		}
		
		m.filteredApps = []AppEntry{webSearchEntry}
		m.offset = 0
	} else {
		if query == "" {
			m.filteredApps = m.apps
		} else {
			m.filteredApps = fuzzyFindApps(query, m.apps)
		}
	}

		if m.cursor >= len(m.filteredApps) {
			m.cursor = 0
			m.offset = 0
		}
		
	case tea.WindowSizeMsg:
	    m.width = msg.Width
	    m.height = msg.Height
	    return m, nil

	}
	return m, cmd
}


func (m model) View() string {
	var b strings.Builder
	b.WriteString(m.textInput.View())
	b.WriteString("\n\n")

	if len(m.filteredApps) == 0 {
		
		// text := b.WriteString("0 results\n")
		text := "0 results\n"
		line := theme.Text.Render(text)
		// line := theme.text.Render(text)
		b.WriteString(line)
	} else {
		endIndex := min(m.offset+10, len(m.filteredApps))

		for i := m.offset; i < endIndex; i++ {
			appName := m.filteredApps[i].Name
			
			var line string

			if m.cursor == i {
				indicator := "> "
				// line = theme.indicator.Render(indicator)
				line = theme.Indicator.Render(indicator)
				text := fmt.Sprintf("%s", appName)
				// line = line + theme.selectedStyle.Render(text)
				line = line + theme.Selected.Render(text)
			} else {
				text := fmt.Sprintf("  %s", appName)
				// line = theme.text.Render(text)
				line = theme.Text.Render(text)
			}
			
			b.WriteString(line + "\n")
		}

		remaining := len(m.filteredApps) - endIndex
		if remaining > 0 {
			var line string
			text := fmt.Sprintf("\n ... and other %d options", remaining)
			line = theme.Text.Render(text)
			b.WriteString(line)
		}
	}

	 return theme.Box.Render(b.String())
}
func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	
	finalModel, err := p.Run()
	if err != nil {
		fmt.Printf("Errore UI: %v", err)
		os.Exit(1)
	}

	if m, ok := finalModel.(model); ok && m.selectedApp != nil {
		launchApp(m.selectedApp)
	}
}

func launchApp(app *AppEntry) {
	cmdStr := app.Exec

	if app.Terminal {
        // Nota: la sintassi flag dipende dal terminale (-e per kitty/alacritty/xterm, -- per gnome-terminal)
		cmdStr = fmt.Sprintf("%s -e %s",  TERMINAL,cmdStr)
	}

	cmd := exec.Command("sh", "-c", cmdStr)
	
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setsid: true,
	}

	cmd.Stdin = nil
	cmd.Stdout = nil
	cmd.Stderr = nil

	err := cmd.Start()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Errore avvio: %v\n", err)
		return
	}

	if err := cmd.Process.Release(); err != nil {
		fmt.Printf("Errore release: %v\n", err)
	}
	
}
