/*
Copyright 2025 The Ketches Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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
