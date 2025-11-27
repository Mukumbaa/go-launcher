package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall" 
	"net/url" 

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)


var TERMINAL string = "alacritty"
var BROWSER string = "firefox"

type model struct {
	textInput    textinput.Model
	apps         []AppEntry
	filteredApps []AppEntry
	cursor       int
	selectedApp  *AppEntry // Campo per memorizzare la scelta
}

func initialModel() model {
	loadConfig()
	ti := textinput.New()
	ti.Placeholder = "Search..."
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 50

	args := os.Args[1:]

	var apps []AppEntry
	if len(args) > 0 {
	    apps = loadFromFile(args[0])
	} else {
	    apps = loadApplications()
	}	// for i:=range apps{
	// 	fmt.Printf("%s (%s) [%s]\n", apps[i].Name, apps[i].File, apps[i].Exec);
	// }
	return model{
		textInput:    ti,
		apps:         apps,
		filteredApps: apps,
		cursor:       0,
		selectedApp:  nil,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit

		case "enter":
			if len(m.filteredApps) > 0 && m.cursor < len(m.filteredApps) {
				// 1. Salva la selezione nel modello
				selected := m.filteredApps[m.cursor]
				m.selectedApp = &selected
				// 2. Esci e basta. Il lancio avverrà nel main.
				return m, tea.Quit
			}

		case "up":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "tab":
			if m.cursor < len(m.filteredApps)-1 {
				m.cursor++
			}

		case "shift+tab":
			if m.cursor > 0 {
				m.cursor--
			}
		}
	}

	m.textInput, cmd = m.textInput.Update(msg)

	rawInput := m.textInput.Value() // Non usare TrimSpace subito per rilevare ":"
	query := strings.TrimSpace(rawInput)

	// Logica: Se inizia con ":", modalità ricerca web
	if strings.HasPrefix(rawInput, "?") {
		searchText := strings.TrimSpace(strings.TrimPrefix(rawInput, "?"))
		
		var displayName string
		var execCmd string
		
		if searchText == "" {
			displayName = "Google search:"
			// Se preme invio ora, apriamo solo il browser vuoto (opzionale)
			execCmd = "google-chrome-stable"			// Se c'è solo ":", svuota la lista o mostra un suggerimento
			m.filteredApps = []AppEntry{} 
		} else {
			// 1. Codifica la query per l'URL (es. spazi diventano + o %20)
			encodedQuery := url.QueryEscape(searchText)
			
			// 2. Costruisci l'URL di Google
			googleUrl := fmt.Sprintf("https://www.google.com/search?q=%s", encodedQuery)
			displayName = fmt.Sprintf("Google search: %s", searchText)
			
			// 3. Crea il comando. 
			// Puoi usare "xdg-open" per il browser predefinito o "google-chrome-stable" per forzare Chrome.
			// Usiamo apici singoli attorno all'URL per sicurezza nella shell.
			execCmd = fmt.Sprintf("%s '%s'", BROWSER,googleUrl)
			
			// 4. Crea una AppEntry "virtuale"
		}
		webSearchEntry := AppEntry{
			Name:     displayName,
			Exec:     execCmd,
			Terminal: false,
			File:     "", // Non esiste un file .desktop reale
		}
		
		// Mostra solo questa opzione
		m.filteredApps = []AppEntry{webSearchEntry}
	} else {
		// Logica originale per le app locali
		if query == "" {
			m.filteredApps = m.apps
		} else {
			m.filteredApps = fuzzyFindApps(query, m.apps)
		}
	}

	// Reset del cursore se la lista cambia drasticamente
	if m.cursor >= len(m.filteredApps) {
		m.cursor = 0
	}

	return m, cmd
}

func (m model) View() string {
	var b strings.Builder
	b.WriteString(m.textInput.View())
	b.WriteString("\n\n")

	if len(m.filteredApps) == 0 {
		b.WriteString("0 results\n")
	} else {
		maxDisplay := min(len(m.filteredApps), 10)

		for i := range maxDisplay {
			cursor := " "
			if m.cursor == i {
				cursor = ">"
			}
			b.WriteString(fmt.Sprintf("%s %s\n", cursor, m.filteredApps[i].Name))
		}
		
		if len(m.filteredApps) > maxDisplay {
			b.WriteString(fmt.Sprintf("\n ... and other %d apps\n", len(m.filteredApps)-maxDisplay))
		}
	}
	return b.String()
}

func main() {
	p := tea.NewProgram(initialModel())
	
	// Esegui il programma e cattura il modello finale
	finalModel, err := p.Run()
	if err != nil {
		fmt.Printf("Errore UI: %v", err)
		os.Exit(1)
	}

	// Controlla se abbiamo selezionato qualcosa
	if m, ok := finalModel.(model); ok && m.selectedApp != nil {
		// Ora siamo fuori da Bubble Tea, il terminale è ripristinato.
		// Possiamo lanciare l'app in sicurezza.
		launchApp(m.selectedApp)
	}
}

// func launchApp(execCmd string) {
// 	// sh -c permette di eseguire comandi con argomenti complessi
	
// 	cmd := exec.Command("sh", "-c", execCmd)
	
// 	// IMPORTANTE: Sgancia il processo dal terminale corrente
// 	cmd.SysProcAttr = &syscall.SysProcAttr{
// 		Setsid: true, // Crea una nuova sessione per il figlio
// 	}

// 	// Rilascia stdin/stdout/stderr per evitare che l'app si blocchi
// 	// aspettando input dal terminale chiuso o scrivendo su un pipe rotto.
// 	cmd.Stdin = nil
// 	cmd.Stdout = nil
// 	cmd.Stderr = nil

// 	err := cmd.Start()
// 	if err != nil {
// 		fmt.Fprintf(os.Stderr, "Errore avvio: %v\n", err)
// 		return
// 	}

// 	// Usiamo Release invece di Wait, così il codice Go può terminare
// 	// senza aspettare che l'app lanciata si chiuda.
// 	// Il processo figlio è ora orfano e adottato da init/systemd.
// 	if err := cmd.Process.Release(); err != nil {
// 		fmt.Printf("Errore release: %v\n", err)
// 	}
	
// 	fmt.Printf("Lanciato: %s\n", execCmd)
// }
// Aggiorna la firma per accettare *AppEntry
func launchApp(app *AppEntry) {
	cmdStr := app.Exec

    // Se l'app richiede un terminale, prepariamo il comando wrapper.
    // Qui dovresti rilevare il terminale dell'utente, hardcodiamo "kitty" per esempio,
    // oppure puoi usare "x-terminal-emulator" -e.
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
	
    // Stampa di debug, puoi rimuoverla
	// fmt.Printf("Lanciato: %s\n", cmdStr) 
}
