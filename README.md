# argparse

In development...

Simple cli argument parser.

## Example

```go
func main() {
    arg := argparse.ArgumentParser()
    arg.SetDescription("Test script")
    arg.SetArgument("arg1", "arg1 info", []string{"a","b"})
    arg.SetArgument("arg2", "arg1 info", []string{})
    arg.Parse()
}
```

Help info example:
```
Usage: /tmp/go-build135913729/b001/exe/sample [--arg1] [--arg2]
Test script

Optional arguments:
--arg1  arg1 info
--arg2  arg2 info
```