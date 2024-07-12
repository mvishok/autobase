package log

import (
	"fmt"

	"github.com/TwiN/go-color"
)

var enabled = true

func Info(message string) {
	if enabled {
		fmt.Println(color.InCyan("INFO: " + message))
	}
}

func Error(message string) {
	if enabled {
		fmt.Println(color.InBold(color.InRed("ERROR: " + message)))
	}
}

func Success(message string) {
	if enabled {
		fmt.Println(color.InBold(color.InGreen("SUCCESS: " + message)))
	}
}

func Warning(message string) {
	if enabled {
		fmt.Println(color.InYellow("WARNING: " + message))
	}
}

func Disable() {
	enabled = false
}

func Enable() {
	enabled = true
}
