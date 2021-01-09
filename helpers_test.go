package argparse

import (
	"testing"
)

func TestHelpers(t *testing.T) {
	arg := makeAlias("arg1")
	if arg != "--arg1" {
		t.Errorf("alias must be --arg1, but got %s", arg)
	}

	str := drawSpaces(5)
	if len(str) != 5 {
		t.Errorf("spaces length must be 5, but got %d", len(str))
	}
}