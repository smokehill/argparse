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
	return &ArgParse{"", "", []*Arg{}}
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
	if len(os.Args) > 1 && len(a.args) > 0 {
		a.parseInput()
	} else {
		a.helpInfo()
	}
}

func (a *ArgParse) helpInfo() {
	name := os.Args[0]
	if a.name != "" {
		name = a.name
	}

	usage := ""
	for _, arg := range a.args {
		alias := alias(arg.name)
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
			alias := alias(arg.name)
			if len(arg.choices) > 0 {
				alias = alias + "=v"
			}
			if len(alias) > maxLen {
				maxLen = len(alias)
			}
		}

		for k, _ := range a.args {
			alias := alias(a.args[k].name)
			help := a.args[k].help
			strLen := 0
			if len(a.args[k].choices) > 0 {
				if a.args[k].choices[0] != "" {
					help = fmt.Sprintf("%s [%s]", help, strings.Join(a.args[k].choices, ","))
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
		alias := alias(arg.name)
		if len(arg.choices) > 0 {
			alias = alias + "=v"
		}
		usage = usage + " [" + alias + "]"
	}

	fmt.Println("Usage: " + name + usage)
	fmt.Println(err)
}

func (a *ArgParse) parseInput() {
	if len(os.Args) > 1 {

		// check arg scheme

		bad := []string{}
		reg1, _ := regexp.Compile("(^--([0-9a-zA-Z-_]+)$|^--([0-9a-zA-Z-_]+)=([0-9a-zA-Z-_]+)$)")

		for _, osArg := range os.Args[1:] {
			if !reg1.MatchString(osArg) {
				bad = append(bad, osArg)
			}
		}

		if len(bad) > 0 {
			a.errorInfo(fmt.Sprintf("Error: bad arguments format %v\n", bad))
			return
		}

		// check arg existence

		reg2 := regexp.MustCompile("=(.*)$")

		for _, osArg := range os.Args[1:] {
			isset := false
			alias1 := reg2.ReplaceAllString(osArg, "")
			for k, _ := range a.args {
				alias2 := alias(a.args[k].name)
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
			a.errorInfo(fmt.Sprintf("Error: unrecognized arguments %v\n", bad))
			return
		}

		// check arg value

		reg3 := regexp.MustCompile("^--([0-9a-zA-Z-_]+)=?")

		for _, osArg := range os.Args[1:] {
			alias1 := reg2.ReplaceAllString(osArg, "")

			for k, _ := range a.args {
				alias2 := alias(a.args[k].name)
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
			a.errorInfo(fmt.Sprintf("Error: bad arguments value %v\n", bad))
			return
		}
	}
}