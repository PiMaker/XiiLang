package interpreter

import (
    "log"
    "github.com/PiMaker/XiiLang/parser"
	"fmt"
	"os"
	"bufio"
    "reflect"
)

func Interpret(nodes []parser.INode, state *parser.XiiState, debug, dump bool) {
    log.Println("Beginning interpretation...")

    if debug {
        fmt.Println("Debugging mode enabled, press enter after each command to continue.")
    }

    log.Println("Initialized state, loop starting now!")

    for {
        tmpNext := state.NextNode.GetID()

        log.Printf("%d/%s\n", tmpNext, reflect.TypeOf(state.NextNode))
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