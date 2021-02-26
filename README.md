# argparse

[![Build Status](https://travis-ci.com/smokehill/argparse.svg?branch=master)](https://travis-ci.com/smokehill/argparse)

Simple cli argument parser.

## Install

```
go get -u github.com/smokehill/argparse
```

## Example

```go
func main() {
    a := argparse.ArgumentParser()
    a.SetName("test")
    a.SetDescription("Test script")
    a.SetArgument("arg1", "arg1 info", []string{"a","b"})   // choices [a,b]
    a.SetArgument("arg2", "arg2 info", []string{""})        // any value
    a.SetArgument("arg3", "arg3 info", []string{})          // no value
    a.Parse()

    // test --arg1=a
    if a.Has("arg1") {
        fmt.Println(a.Get("arg1")) // a
    }

    // test --arg2=123
    if a.Has("arg2") {
        fmt.Println(a.Get("arg2")) // 123
    }

    // test --arg3
    if a.Has("arg3") {
        fmt.Println(a.Get("arg3")) // ""
    }
}
```

Help output:
```
Usage: test [--help] [--arg1=v] [--arg2=v] [--arg3]
Test script

Optional arguments:
--help   show this help message
--arg1=v info arg1 [a,b]
--arg2=v info arg2
--arg3   info arg3
```