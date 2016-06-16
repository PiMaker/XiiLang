package interpreter

import (
    "log"
    "github.com/PiMaker/XiiLang/parser"
	"fmt"
	"os"
	"bufio"
    "reflect"
    "time"
)

func Interpret(nodes []parser.INode, state *parser.XiiState, debug, dump, trace, time bool) {
    log.Println("Beginning interpretation...")

    if debug || dump || trace || time {
        log.Println("Using debug interpreter, expect performance penalties.")
        InterpretDebug(nodes, state, debug, dump, trace, time)
    } else {
        log.Println("Using release interpreter.")
        InterpretRelease(nodes, state)
    }
}

func InterpretRelease(nodes []parser.INode, state *parser.XiiState) {
    for {
        tmpNext := state.NextNode.GetID()

        err := state.NextNode.Execute(state)
        
        if err != nil {
            fmt.Println("Error: " + err.Error())
            break
        }

        if state.NextNode != nil && tmpNext == state.NextNode.GetID() {
            state.NextNode = state.NextNode.Next()
        }

        if state.NextNode == nil {
            break
        }
    }
}

var ExecutionTimeTable map[string][]time.Duration

func InterpretDebug(nodes []parser.INode, state *parser.XiiState, debug, dump, trace, timeExec bool) {
    if debug {
        fmt.Println("Debugging mode enabled, press enter after each command to continue.")
    }

    log.Println("Initialized state, loop starting now!")

    if trace {
        fmt.Println("Trace enabled")
    }

    if timeExec {
        ExecutionTimeTable = make(map[string][]time.Duration)
    }

    for {
        tmpNext := state.NextNode.GetID()

        if trace {
            fmt.Printf("#Trace :: ID: %d / %s / %s\n", tmpNext, state.NextNode.GetTrace(), reflect.TypeOf(state.NextNode))
        }

        var beforeTime time.Time
        if timeExec {
            beforeTime = time.Now()
        }

        err := state.NextNode.Execute(state)

        if timeExec {
            duration := time.Since(beforeTime)
            keyword := state.NextNode.GetTrace()
            ExecutionTimeTable[keyword] = append(ExecutionTimeTable[keyword], duration)
        }
        
        if err != nil {
            fmt.Println("Error: " + err.Error())
            break
        }

        if state.NextNode != nil && tmpNext == state.NextNode.GetID() {
            state.NextNode = state.NextNode.Next()
        }

        if state.NextNode == nil {
            break
        }

        if dump {
            fmt.Println()
            fmt.Println("---- Literal-Table:")
            dumpStringMap(state.VariableLiteralTable)
            fmt.Println("---- Variable-Table:")
            dumpFloatMap(state.VariableNumberTable)
        }

        if debug {
            fmt.Println("Command done.")
            reader := bufio.NewReader(os.Stdin)
            _, _ = reader.ReadString('\n')
            fmt.Println("Starting new command...")
        }
    }
}

func dumpStringMap(m map[string]string) {
    for k, v := range m {
        fmt.Printf("%s:\t%s\n", k, v)
    }
    fmt.Println()
}

func dumpFloatMap(m map[string]float64) {
    for k, v := range m {
        fmt.Printf("%s:\t%f\n", k, v)
    }
    fmt.Println()
}