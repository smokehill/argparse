package argparse

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

// ArgParse argument parser instance.
type ArgParse struct {
	// Program name (optional).
	name string
	// Program description (optional).
	description string
	// Arguments list.
	args []*Arg
}

// Arg argument instance.
type Arg struct {
	// Argument name.
	name string
	// Argument description (optional).
	help string
	// Argument choices (optional).
	choices []string
	// Argument value from input.
	value string
	// True if input argument exists.
	active bool
}

// ArgumentParser return ArgParse instance.
func ArgumentParser() *ArgParse {
	return &ArgParse{"", "", []*Arg{}}
}

// SetName sets program name.
func (a *ArgParse) SetName(name string) *ArgParse {
	a.name = name
	return a
}

// SetDescription sets program description.
func (a *ArgParse) SetDescription(description string) *ArgParse {
	a.description = description
	return a
}

// SetArgument defines a new argument.
func (a *ArgParse) SetArgument(name string, help string, choices []string) *ArgParse {
	a.args = append(a.args, &Arg{name, help, choices, "", false})
	return a
}

// Has checks if input argument is preset.
func (a *ArgParse) Has(name string) bool {
	isset := false
	for _, arg := range a.args {
		if arg.name == name && arg.active {
			isset = true
			break
		}
	}
	return isset
}

// Get returns argument value.
func (a *ArgParse) Get(name string) string {
	value := ""
	for _, arg := range a.args {
		if arg.name == name {
			value = arg.value
			break
		}
	}
	return value
}

// Parse handles input and parses available arguments data.
func (a *ArgParse) Parse() {
	if len(a.args) > 0 {
		a.checkArgName()
		a.checkArgChoices()
	}
	if len(os.Args) > 1 {
		a.parseInput()
	} else {
		a.helpInfo()
	}
}

// helpInfo displays help information.
func (a *ArgParse) helpInfo() {
	name := os.Args[0]
	if a.name != "" {
		name = a.name
	}

	usage := ""
	for _, arg := range a.args {
		alias := makeAlias(arg.name)
		if len(arg.choices) > 0 {
			alias = alias + "=v"
		}
		usage = usage + " [" + alias + "]"
	}

	fmt.Println("Usage: " + name + usage)

	if a.description != "" {
		fmt.Println(a.description)
	}

	fmt.Println("")
	fmt.Println("Optional arguments:")

	if len(a.args) > 0 {
		maxLen := 0
		for _, arg := range a.args {
			alias := makeAlias(arg.name)
			if len(arg.choices) > 0 {
				alias = alias + "=v"
			}
			if len(alias) > maxLen {
				maxLen = len(alias)
			}
		}

		for k, _ := range a.args {
			alias := makeAlias(a.args[k].name)
			help := a.args[k].help
			strLen := 0
			if len(a.args[k].choices) > 1 {
				if help != "" {
					help = help + " "
				}
				help = fmt.Sprintf("%s[%s]", help, strings.Join(a.args[k].choices, ","))
				alias = alias + "=v"
			}
			if len(alias) < maxLen {
				strLen = maxLen - len(alias)
			}

			fmt.Printf(alias + drawSpaces(strLen) + " %s\n", help)
		}
	}

	fmt.Println("")
}

// errorInfo displays error information.
func (a *ArgParse) errorInfo(err string) {
	name := os.Args[0]
	if a.name != "" {
		name = a.name
	}

	usage := ""
	for _, arg := range a.args {
		alias := makeAlias(arg.name)
		if len(arg.choices) > 0 {
			alias = alias + "=v"
		}
		usage = usage + " [" + alias + "]"
	}

	fmt.Println("Usage: " + name + usage)
	fmt.Println(err)
	fmt.Println("")

	os.Exit(0)
}

// checkArgName checks argument name.
func (a *ArgParse) checkArgName() {
	bad := []string{}
	reg, _ := regexp.Compile("^[0-9a-zA-Z-_]+$")

	for _, arg := range a.args {
		if !reg.MatchString(arg.name) {
			alias := makeAlias(arg.name)
			bad = append(bad, alias)
		}
	}

	if len(bad) > 0 {
		err := "Error: bad arguments names: " + strings.Join(bad, ", ")
		a.errorInfo(err)
	}
}

// checkArgChoices checks argument choices.
func (a *ArgParse) checkArgChoices() {
	bad := []string{}

	for _, arg := range a.args {
		if len(arg.choices) == 1 && arg.choices[0] != "" {
			bad = append(bad, makeAlias(arg.name) + "=v")
		} else if len(arg.choices) > 1 {
			for _, v := range arg.choices {
				if v == "" {
					bad = append(bad, fmt.Sprintf("%s=v %v", makeAlias(arg.name), arg.choices))
				}
			}
		}
	}

	if len(bad) > 0 {
		err := "Error: bad arguments choices: " + strings.Join(bad, ", ")
		a.errorInfo(err)
	}
}

// parseInput handles input.
func (a *ArgParse) parseInput() {
	bad := []string{}

	// Check input argument scheme.

	reg1, _ := regexp.Compile("(^--([0-9a-zA-Z-_]+)$|^--([0-9a-zA-Z-_]+)=([0-9a-zA-Z-_]+)$)")

	for _, osArg := range os.Args[1:] {
		if !reg1.MatchString(osArg) {
			bad = append(bad, osArg)
		}
	}

	if len(bad) > 0 {
		err := "Error: bad arguments format: " + strings.Join(bad, ", ")
		a.errorInfo(err)
	}

	// Recognize input argument.

	reg2 := regexp.MustCompile("=(.*)$")

	for _, osArg := range os.Args[1:] {
		isset := false
		alias1 := reg2.ReplaceAllString(osArg, "")

		for k, _ := range a.args {
			alias2 := makeAlias(a.args[k].name)

			if alias1 == alias2 {
				a.args[k].active = true
				isset = true
				break
			}
		}
		if !isset {
			bad = append(bad, osArg)
		}
	}

	if len(bad) > 0 {
		err := "Error: unrecognized arguments: " + strings.Join(bad, ", ")
		a.errorInfo(err)
	}

	// Check input argument value.

	reg3 := regexp.MustCompile("^--([0-9a-zA-Z-_]+)=?")

	for _, osArg := range os.Args[1:] {
		alias1 := reg2.ReplaceAllString(osArg, "")

		for k, _ := range a.args {
			alias2 := makeAlias(a.args[k].name)
			value := reg3.ReplaceAllString(osArg, "")

			if alias1 == alias2 {
				if len(a.args[k].choices) > 0 {
					if value == "" {
						bad = append(bad, osArg)
						break
					}

					if a.args[k].choices[0] != "" {
						isset := false
						for _, val := range a.args[k].choices {
							if val == value {
								isset = true
							}
						}

						if !isset {
							bad = append(bad, osArg)
						}
					}

					a.args[k].value = value
				} else {
					if value != "" {
						bad = append(bad, osArg)
					}
				}
			}
		}
	}

	if len(bad) > 0 {
		err := "Error: bad arguments value: " + strings.Join(bad, ", ")
		a.errorInfo(err)
	}
}