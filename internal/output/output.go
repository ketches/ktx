package output

import (
	"os"

	"github.com/fatih/color"
)

// Fatal prints a fatal error message and exits the program.
func Fatal(format string, a ...any) {
	color.Red("😾 "+format, a...)
	os.Exit(1)
}

// Note prints a note message.
func Note(format string, a ...interface{}) {
	color.New(color.Faint).Printf("😼 "+format+"\n", a...)
}

// Done prints a done message.
func Done(format string, a ...interface{}) {
	color.Green("😺 "+format, a...)
}

// Fail prints a fail message.
func Fail(format string, a ...interface{}) {
	color.Red("😾 "+format, a...)
}
