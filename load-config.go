package main


import (
	"bufio"
	"os"
	"strings"
)

var pathConfig = "/home/mukumba/.config/go-launcher/go-launcher.config"
func loadConfig(){
	
	file, err := os.Open(pathConfig)
	check(err)
	
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan(){

		line := scanner.Text()
			
		if strings.HasPrefix(line, "Terminal=") {
			TERMINAL = line[9:]
		} else if strings.HasPrefix(line, "Browser="){
			BROWSER = line[8:]
		}
	}
}
