package cdtacctwebscraper

import (
	"fmt"
	"strings"
)

func GenerateCommands(inlist [][]string) (commands Commands, err error) {

	// Initialize stack
	root := Command{}
	stack := []*Command{&root}
	index := 0

	for _, list := range inlist {
		n := len(list)
		if n == 0 {
			err = fmt.Errorf("Invalid command %q", list)
			return
		}

		// Record and remove indent
		full := len(list[0])
		command := Command{Instruction: strings.TrimLeft(list[0], " "), Params: list[1:]}
		indent := full - len(command.Instruction)

		if indent%4 != 0 {
			err = fmt.Errorf("Unexpected indent %d chars for %q", indent, list)
			return
		}
		level := indent / 4

		// Push subroutine
		if level > index {
			if level > index+1 {
				err = fmt.Errorf("Skipped forward to subroutine level %d for %q", level, list)
				return
			}
			l := len(stack[index].Commands)
			if l == 0 {
				err = fmt.Errorf("Can't push subroutine level %d without previous command for %q", level, list)
				return
			}
			stack = append(stack, &stack[index].Commands[l-1])
			index = level
		}

		// Pop back
		if level < index {
			if level < index-1 {
				err = fmt.Errorf("Skipped back to subroutine level %d for %q", level, list)
				return
			}
			l := len(stack)
			stack = stack[:l-1]
			index = level
		}

		stack[index].Commands = append(stack[index].Commands, command)
	}

	commands = root.Commands
	return
}

func (browser *Browser) PrintCommands(commands Commands, indent int) {
	for _, command := range commands {
		pref := strings.Repeat("    ", indent)
		browser.Log.Printf("%s  %s(%q)\n", pref, command.Instruction, command.Params)
		if command.Commands != nil {
			browser.PrintCommands(command.Commands, indent+1)
		}
	}
}
