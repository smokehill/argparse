package argparse

func makeAlias(val string) string {
	return "--" + val
}

func drawSpaces(num int) string {
	str := ""
	for num > 0 {
		str = str + " "
		num--
	}
	return str
}