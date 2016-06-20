package interpreter

import (
    "log"
	"fmt"
	"os"
	"bufio"
    "reflect"
    "time"
)

func Interpret(nodes []INode, state *XiiState, debug, trace, time bool) {
    log.Println("Beginning interpretation...")

    if debug || trace || time {
        log.Println("Using debug interpreter, expect performance penalties.")
        InterpretDebug(nodes, state, debug, trace, time)
    } else {
        log.Println("Using release interpreter.")
        InterpretRelease(nodes, state)
    }
}

func InterpretRelease(nodes []INode, state *XiiState) {
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

func InterpretDebug(nodes []INode, state *XiiState, debug, trace, timeExec bool) {
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
            fmt.Printf("Trace :: ID: %d / %s / %s\n", tmpNext, state.NextNode.GetTrace(), reflect.TypeOf(state.NextNode))
        }

        var beforeTime time.Time
        if timeExec {
            beforeTime = time.Now()
        }

        err := state.NextNode.Execute(state)

        if timeExec && state.NextNode != nil {
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

        if debug {
            fmt.Println("Command done.")
            reader := bufio.NewReader(os.Stdin)
            _, _ = reader.ReadString('\n')
            fmt.Println("Starting new command...")
        }
    }
}

func dumpMap(m map[string]interface{}) {
    for k, v := range m {
        fmt.Printf("%s:\t%s\n", k, v)
    }
    fmt.Println()
}