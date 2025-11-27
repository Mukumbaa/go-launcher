package main

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type AppEntry struct {
	Name     string
	Exec     string
	File     string
	Terminal bool
}

func loadApplications() []AppEntry {
	paths := []string{
		filepath.Join(os.Getenv("HOME"), ".local/share/applications"),
		"/usr/local/share/applications",
		"/usr/share/applications",
	}
	
	var (
		mu   sync.Mutex
		wg   sync.WaitGroup
		apps []AppEntry
		seen = make(map[string]bool)
	)
	
	semaphore := make(chan struct{}, 20)
	
	for _, dir := range paths {
		entries, err := os.ReadDir(dir)
		if err != nil {
			continue
		}
		
		for _, entry := range entries {
			if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".desktop") {
				continue
			}
			
			mu.Lock()
			if seen[entry.Name()] {
				mu.Unlock()
				continue
			}
			seen[entry.Name()] = true
			mu.Unlock()
			
			wg.Add(1)
			go func(dir string, name string) {
				defer wg.Done()
				
				semaphore <- struct{}{}
				defer func() { <-semaphore }()
				
				path := filepath.Join(dir, name)
				if app, ok := parseDesktopFile(path); ok {
					mu.Lock()
					apps = append(apps, app)
					mu.Unlock()
				}
			}(dir, entry.Name())
		}
	}
	
	wg.Wait()
	return apps
}

func parseDesktopFile(path string) (AppEntry, bool) {
	f, err := os.Open(path)
	if err != nil {
		return AppEntry{}, false
	}
	defer f.Close()
	
	var (
		inDesktopEntry bool
		name           string
		exec           string
		noDisplay      bool
		hidden         bool
		isApp          bool
		terminal       bool
		onlyShowIn     string
		notShowIn      string
	)
	
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		
		if len(line) == 0 || line[0] == '#' {
			continue
		}
		
		if line[0] == '[' {
			inDesktopEntry = line == "[Desktop Entry]"
			continue
		}
		
		if !inDesktopEntry {
			continue
		}
		
		if noDisplay || hidden {
			break
		}
		
		if strings.HasPrefix(line, "Type=") {
			isApp = line[5:] == "Application"
			if !isApp {
				return AppEntry{}, false
			}
		} else if strings.HasPrefix(line, "NoDisplay=") {
			noDisplay = line[10:] == "true"
			if noDisplay {
				return AppEntry{}, false
			}
		} else if strings.HasPrefix(line, "Hidden=") {
			hidden = line[7:] == "true"
			if hidden {
				return AppEntry{}, false
			}
		} else if strings.HasPrefix(line, "OnlyShowIn=") {
			onlyShowIn = line[11:]
		} else if strings.HasPrefix(line, "NotShowIn=") {
			notShowIn = line[10:]
		} else if strings.HasPrefix(line, "Terminal=") {
			terminal = line[9:] == "true"
		} else if name == "" && strings.HasPrefix(line, "Name=") {
			name = line[5:]
		} else if exec == "" && strings.HasPrefix(line, "Exec=") {
			exec = line[5:]
		}
	}
	
	if !isApp || name == "" || exec == "" {
		return AppEntry{}, false
	}
	
	if !isCompatibleWithHyprland(onlyShowIn, notShowIn) {
		return AppEntry{}, false
	}
	
	exec = removeExecPlaceholders(exec)
	
	return AppEntry{
		Name:     name,
		Exec:     exec,
		File:     path,
		Terminal: terminal,
	}, true
}

func isCompatibleWithHyprland(onlyShowIn, notShowIn string) bool {
	if notShowIn != "" {
		desktops := strings.SplitSeq(notShowIn, ";")
		for de := range desktops {
			de = strings.TrimSpace(de)
			if de == "Hyprland" || de == "wlroots" {
				return false
			}
		}
	}
	
	if onlyShowIn != "" {
		desktops := strings.SplitSeq(onlyShowIn, ";")
		
		for de := range desktops {
			de = strings.TrimSpace(de)
			// Se trova Hyprland, wlroots o X-Generic, è compatibile
			if de == "Hyprland" || de == "wlroots" || de == "X-Generic" {
				return true
			}
		}
		
		// Se OnlyShowIn è specificato ma non include Hyprland/wlroots/X-Generic,
		// allora NON è compatibile (esclude GNOME, KDE, Unity, ecc.)
		return false
	}
	
	return true
}

func removeExecPlaceholders(exec string) string {
	var result strings.Builder
	result.Grow(len(exec))
	
	inWord := false
	skipWord := false
	
	for i := 0; i < len(exec); i++ {
		c := exec[i]
		
		if c == ' ' || c == '\t' {
			if inWord && !skipWord {
				result.WriteByte(' ')
			}
			inWord = false
			skipWord = false
		} else {
			if !inWord {
				inWord = true
				skipWord = (c == '%')
			}
			if !skipWord {
				result.WriteByte(c)
			}
		}
	}
	
	return strings.TrimSpace(result.String())
}

func check(e error) {
    if e != nil {
        panic(e)
    }
}



func loadFromFile(path string) []AppEntry{
	file, err := os.Open(path)
	check(err)

	defer file.Close()

	scanner := bufio.NewScanner(file)

	apps := []AppEntry{}

	for scanner.Scan() {
		var (
			name           string
			exec           string
			terminal       bool
		)
		line := scanner.Text()
		app := strings.SplitSeq(line, "|")

		for op := range app{
			
			if strings.HasPrefix(op, "Terminal=") {
				terminal = op[9:] == "true"
			} else if name == "" && strings.HasPrefix(op, "Name=") {
				name = op[5:]
			} else if exec == "" && strings.HasPrefix(op, "Exec=") {
				exec = op[5:]
			}
		}
			
		apps = append(apps, AppEntry{
			Name: name,
			Exec: exec,
			Terminal: terminal,	
		})
	}
	// fmt.Printf("PROVA\n");
	// for i := 0; i<len(apps); i++{
	// 	fmt.Printf("%s|%s|%t\n", apps[i].Name, apps[i].Exec, apps[i].Terminal)
	// }
	return apps
}

