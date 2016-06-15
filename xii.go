package main

import (
    "io/ioutil"
    "fmt"
    "flag"
    "log"
    "github.com/PiMaker/XiiLang/parser"
    "github.com/PiMaker/XiiLang/interpreter"
)

func main() {
    fmt.Println("XiiLang(sr) v0.1, (C) Stefan Reiter 2016")

    verbose := flag.Bool("v", false, "Be verbose with output")
    debug := flag.Bool("d", false, "Use the debug mode to execute a script line by line and wait for enter")
    dump := flag.Bool("p", false, "Print table information after each node")

    flag.Parse()

    verboseVal := *verbose
    path := flag.Arg(0)

    if !verboseVal {
        log.SetOutput(ioutil.Discard)
    }

    log.Println("Loading file: " + path)

    tokens, err := parser.TokenizeFile(path)

    if err != nil {
        fmt.Println(err.Error())
        return
    }

    nodes, err := parser.ParseTokens(tokens)
    if err != nil {
        fmt.Println(err.Error())
        return
    }

    log.Printf("AST: %s\n", nodes)

    fmt.Println()

    state := &parser.XiiState{}
    state.VariableNumberTable = make(map[string]float64)
    state.VariableLiteralTable = make(map[string]string)
    state.Nodes = nodes
    state.NextNode = nodes[0]

    interpreter.Interpret(nodes, state, *debug, *dump)

    fmt.Println()
    fmt.Println("XiiLang: Execution ended")
}