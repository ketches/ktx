package output

import (
	"os"

	"github.com/fatih/color"
)

func Fatal(format string, a ...any) {
	color.Red("😾 "+format, a...)
	os.Exit(1)
}

func Note(format string, a ...interface{}) {
	color.New(color.Faint).Printf("😼 "+format+"\n", a...)
}

func Done(format string, a ...interface{}) {
	color.Green("😺 "+format, a...)
}

func Fail(format string, a ...interface{}) {
	color.Red("😾 "+format, a...)
}
