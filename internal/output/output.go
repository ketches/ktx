package output

import (
	"os"

	"github.com/fatih/color"
)

func Fatal(s string, a ...any) {
	color.Red("😹 "+s, a...)
	os.Exit(1)
}

func Note(s string, a ...interface{}) {
	color.New(color.Faint).Printf("😼 "+s+"\n", a...)
}

func Done(s string, a ...interface{}) {
	color.Green("😺 "+s, a...)
}

func Fail(s string, a ...interface{}) {
	color.Red("😾 "+s, a...)
}
