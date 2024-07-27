package log

import (
	"fmt"
)

var enabled = true

const Reset = "\033[0m"
const Red = "\033[31m"
const Green = "\033[32m"
const Yellow = "\033[33m"
const Blue = "\033[34m"
const Magenta = "\033[35m"
const Cyan = "\033[36m"
const Gray = "\033[37m"
const White = "\033[97m"

func Info(message string) {
	if enabled {
		fmt.Println(Cyan + "INFO: " + Reset + message)
	}
}

func Error(message string) {
	if enabled {
		fmt.Println(Red + "ERROR: " + Reset + message)
	}
}

func Success(message string) {
	if enabled {
		fmt.Println(Green + "SUCCESS: " + Reset + message)
	}
}

func Warning(message string) {
	if enabled {
		fmt.Println(Yellow + "WARNING: " + Reset + message)
	}
}

func Disable() {
	enabled = false
}

func Enable() {
	enabled = true
}
