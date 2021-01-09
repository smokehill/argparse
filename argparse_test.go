package argparse

import (
	"os"
	"testing"
)

func TestArgumentParser(t *testing.T) {
	os.Args = []string{"test","--arg1=123","--arg2"}

	a := ArgumentParser()
	a.SetName("test")
	a.SetDescription("Test script")
	a.SetArgument("arg1", "info arg2", []string{""})
	a.SetArgument("arg2", "info arg3", []string{})
	a.Parse()

	if !a.Has("arg1") {
		t.Errorf("a.Has(\"arg1\") must be 'true', but got 'false'")
	}

	if a.Get("arg1") != "123" {
		t.Errorf("a.Get(\"arg1\") must be '123', but got '%s'", a.Get("arg1"))
	}

	if !a.Has("arg2") {
		t.Errorf("a.Has(\"arg2\") must be 'true', but got 'false'")
	}
}