package log

import "fmt"

const blue = "\033[34m"
const red = "\033[31m"
const green = "\033[32m"
const yellow = "\033[33m"
const reset = "\033[0m"

var enabled = true

func Info(message string) {
	if enabled {
		fmt.Println(blue + "INFO:" + message + reset)
	}
}

func Error(message string) {
	if enabled {
		fmt.Println(red + "ERROR:" + message + reset)
	}
}

func Success(message string) {
	if enabled {
		fmt.Println(green + "SUCCESS:" + message + reset)
	}
}

func Warning(message string) {
	if enabled {
		fmt.Println(yellow + "WARNING:" + message + reset)
	}
}

func Disable() {
	enabled = false
}

func Enable() {
	enabled = true
}
