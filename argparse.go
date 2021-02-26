package argparse

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

type ArgParse struct {
	name string
	description string
	args []*Arg
}

type Arg struct {
	name string
	help string
	choices []string
	value string
	active bool
}

func ArgumentParser() *ArgParse {
	args := []*Arg{}
	args = append(args, &Arg{"help", "show this help message", []string{}, "", false})
	return &ArgParse{"", "", args}
}

func (a *ArgParse) SetName(name string) *ArgParse {
	a.name = name
	return a
}

func (a *ArgParse) SetDescription(description string) *ArgParse {
	a.description = description
	return a
}

func (a *ArgParse) SetArgument(name string, help string, choices []string) *ArgParse {
	a.args = append(a.args, &Arg{name, help, choices, "", false})
	return a
}

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

func (a *ArgParse) Parse() {
	if len(a.args) > 0 {
		a.checkArgName()
		a.checkArgChoices()
	}
	if len(os.Args) > 1 {
		a.parseInput()
		// show help menu
		if os.Args[1] == "--help" {
			a.helpInfo()
		}
	} else {
		name := os.Args[0]
		if a.name != "" {
			name = a.name
		}
		fmt.Println("Use: " + name + " --help")
	}
}

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
			if len(a.args[k].choices) > 0 {
				if len(a.args[k].choices) > 1 {
					if help != "" {
						help = help + " "
					}
					help = fmt.Sprintf("%s[%s]", help, strings.Join(a.args[k].choices, ","))
				}
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

	os.Exit(1)
}

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

func (a *ArgParse) parseInput() {
	bad := []string{}

	// check input argument scheme
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

	// recognize input argument
	reg2 := regexp.MustCompile("=(.*)$")
	for _, osArg := range os.Args[1:] {
		isset := false
		alias1 := reg2.ReplaceAllString(osArg, "")
		for k, _ := range a.args {
			alias2 := makeAlias(a.args[k].name)
			if alias1 == alias2 {
				// make argument active
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

	// check input argument value (any value and static choices)
	reg3 := regexp.MustCompile("^--([0-9a-zA-Z-_]+)=?")
	for _, osArg := range os.Args[1:] {
		alias1 := reg2.ReplaceAllString(osArg, "")
		for k, _ := range a.args {
			alias2 := makeAlias(a.args[k].name)
			value := reg3.ReplaceAllString(osArg, "")
			if alias1 == alias2 {
				if len(a.args[k].choices) > 0 {
					// check empty value
					if value == "" {
						bad = append(bad, osArg)
						break
					}
					// check value in choices
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
					// remember argument value
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