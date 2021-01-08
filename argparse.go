package argparse

import (
	"fmt"
	"os"
	"regexp"
)

type ArgParse struct {
	name string
	args map[string] *Arg
	description string
}

type Arg struct {
	name string
	help string
	choices []string
	value string
	active bool
}

func ArgumentParser() *ArgParse {
	return &ArgParse{"", map[string] *Arg{}, ""}
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
	a.args[alias(name)] = &Arg{name, help, choices, "", false}
	return a
}

func (a *ArgParse) Has(name string) bool {
	isset := false
	if _, ok := a.args[alias(name)]; ok {
		if a.args[alias(name)].active {
			isset = true
		}
	}
	return isset
}

func (a *ArgParse) Get(name string) string {
	value := ""
	if a.Has(name) {
		value = a.args[alias(name)].value
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
	for alias, _ := range a.args {
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
		for alias, _ := range a.args {
			if len(alias) > maxLen {
				maxLen = len(alias)
			}
		}

		for alias, arg := range a.args {
			sLen := 0
			if len(alias) < maxLen {
				sLen = maxLen - len(alias)
			}

			fmt.Printf(alias + drawSpaces(sLen) + " %s\n", arg.help)
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
	for alias, _ := range a.args {
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

		bad = []string{}
		reg2 := regexp.MustCompile("=(.*)$")

		for _, osArg := range os.Args[1:] {
			isset := false
			alias := reg2.ReplaceAllString(osArg, "")
			if _, ok := a.args[alias]; ok {
				isset = true
				a.args[alias].active = true
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

		bad = []string{}
		reg3, _ := regexp.Compile("^--([0-9a-zA-Z-_]+)=([0-9a-zA-Z-_]+)$")
		reg4 := regexp.MustCompile("^--([0-9a-zA-Z-_]+)=")

		for _, osArg := range os.Args[1:] {
			alias := reg2.ReplaceAllString(osArg, "")

			if reg3.MatchString(osArg) {
				value := reg4.ReplaceAllString(osArg, "")
				a.args[alias].value = value

				if len(a.args[alias].choices) > 0 {
					isset := false

					for _, v := range a.args[alias].choices {
						if v == value {
							isset = true
						}
					}

					if !isset {
						bad = append(bad, osArg)
					}
				}
			} else {
				if len(a.args[alias].choices) > 0 {
					bad = append(bad, osArg)
				}
			}
		}

		if len(bad) > 0 {
			a.errorInfo(fmt.Sprintf("Error: bad arguments value %v\n", bad))
			return
		}
	}
}