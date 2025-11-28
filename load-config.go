package main


import (
	"bufio"
	"os"
	"strings"
	"path/filepath"
)

func loadConfig(){
	
	configDir, err := os.UserConfigDir()
	if err != nil {
		return // default if no file
	}

	//path build
	path := filepath.Join(configDir, "go-launcher", "go-launcher.conf")

	file, err := os.Open(path)
	if err != nil {
		// default if no file
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan(){

		line := scanner.Text()

		// ignore empty string or comment
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch key{
			case "Browser":
				BROWSER = value
			case "Terminal":
				TERMINAL = value
			case "Theme":
				THEME = value
		}



		// if strings.HasPrefix(line, "Terminal=") {
		// 	TERMINAL = line[9:]
		// } else if strings.HasPrefix(line, "Browser="){
		// 	BROWSER = line[8:]
		// } else if strings.HasPrefix(line, "Theme="){
		// 	THEME = line[6:]
		// }
	}
}
