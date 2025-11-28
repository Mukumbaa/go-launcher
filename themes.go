package main

import (
	"bufio"
	"os"

	"github.com/charmbracelet/lipgloss"
)

type SelectedStyle struct{
	backgroundColor string
	textColor       string
	indicator       string
}
type TextStyle struct{
	textColor string
}
type InputStyle struct{
	textColor        lipgloss.Style
	placeholderColor lipgloss.Style
	promptColor      lipgloss.Style
}
type Theme struct{
	boxStyle      lipgloss.Style
	text          lipgloss.Style 
	selectedStyle lipgloss.Style
	inputStyle    InputStyle
	indicator     lipgloss.Style
}



var selectedStyle = SelectedStyle{}
var text = TextStyle{}
var input = InputStyle{}

func getDefault(){
	
	selectedStyle.backgroundColor = "#ffffff"
	selectedStyle.textColor = "#000000"
	text.textColor = "#ffffff"
	input.textColor = lipgloss.NewStyle().Foreground(lipgloss.Color("#ffffff"))
	input.placeholderColor = lipgloss.NewStyle().Foreground(lipgloss.Color("#808080"))
	input.promptColor = lipgloss.NewStyle().Foreground(lipgloss.Color("#ffffff"))
	selectedStyle.indicator = "#ffffff"
}

func loadTheme(themeName string) Theme{

	switch themeName{
		case "default":
			// selectedStyle.backgroundColor = "#ffffff"
			// selectedStyle.textColor = "#000000"
			// text.textColor = "#ffffff"
			// input.textColor = lipgloss.NewStyle().Foreground(lipgloss.Color("#ffffff"))
			// input.placeholderColor = lipgloss.NewStyle().Foreground(lipgloss.Color("#808080"))
			// input.promptColor = lipgloss.NewStyle().Foreground(lipgloss.Color("#ffffff"))
			// selectedStyle.indicator = "#ffffff"
			getDefault()

		case "rose-pine":
			selectedStyle.backgroundColor = "#403d52"
			selectedStyle.textColor = "#e0def4"
			// text.textColor = "#e0def4"
			text.textColor = "#908caa"
			input.textColor = lipgloss.NewStyle().Foreground(lipgloss.Color("#e0def4"))
			input.placeholderColor = lipgloss.NewStyle().Foreground(lipgloss.Color("#808080"))
			input.promptColor = lipgloss.NewStyle().Foreground(lipgloss.Color("#c4a7e7"))
			selectedStyle.indicator = "#c4a7e7"
			
	}

	
	return Theme{
		boxStyle: lipgloss.NewStyle().
			Padding(1, 2).
			Margin(1, 0),
		text: lipgloss.NewStyle().
			Foreground(lipgloss.Color(text.textColor)),
		selectedStyle: lipgloss.NewStyle().
			Background(lipgloss.Color(selectedStyle.backgroundColor)).
			Foreground(lipgloss.Color(selectedStyle.textColor)).
			Bold(true),
		inputStyle: input,
		indicator: lipgloss.NewStyle().Foreground(lipgloss.Color(selectedStyle.indicator)),
	}
}

func loadThemeFromFile(){
	file, err := os.Open("/home/mukumba/.config/go-launcher/theme.conf")
	if err != nil {
        return
    }

	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan(){
		
	}
    
}
