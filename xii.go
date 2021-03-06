package main

import (
    "io/ioutil"
    "os"
    "bufio"
    "fmt"
    "flag"
    "time"
    "log"
    "sort"
    "github.com/PiMaker/XiiLang/interpreter"
)

func main() {
    fmt.Println("XiiLang(sr) v0.5, (C) Stefan Reiter 2016\n")

    verbose := flag.Bool("v", false, "Be verbose with output")
    debug := flag.Bool("d", false, "Execute a script line by line and wait for enter")
    trace := flag.Bool("t", false, "Trace mode, prints statement information for every executed node")
    verboseEval := flag.Bool("e", false, "Trace eval calls for conditions")
    stats := flag.Bool("s", false, "Print runtime stats after execution")

    flag.Parse()

    verboseVal := *verbose
    path := flag.Arg(0)

    if !verboseVal {
        log.SetOutput(ioutil.Discard)
    }

    interpreter.VerboseEval = *verboseEval

    if *verboseEval {
        fmt.Println("Eval trace enabled")
    }

    log.Println("Loading file: " + path)

    startTime := time.Now()

    tokens, err := interpreter.TokenizeFile(path)

    if err != nil {
        fmt.Println(err.Error())
        return
    }

    nodes, err := interpreter.ParseTokens(tokens)
    if err != nil {
        fmt.Println(err.Error())
        return
    }

    log.Printf("AST: %s\n", nodes)

    log.Printf("Compilation took %s\n", time.Since(startTime))

    startTime = time.Now()

    state := &interpreter.XiiState{}
    state.Nodes = nodes
    state.NextNode = nodes[0]
    state.FunctionStack = interpreter.NewNodeStack()
    state.StdOut = bufio.NewWriter(os.Stdout)

    interpreter.Interpret(nodes, state, *debug, *trace, *stats)

    state.StdOut.Flush()

    fmt.Println()
    fmt.Println("XiiLang: Execution ended")

    if !*stats {
        log.Printf("Execution time: %s\n", time.Since(startTime))
    }

    if *stats {
        fmt.Println()
        fmt.Println("Execution stats:")

        fmt.Printf("- Execution time total: %s\n", time.Since(startTime))

        slice := convertToSortedSlice(interpreter.ExecutionTimeTable)
        for _, v := range slice {
            fmt.Printf("- %s: \tTotal: %s / Avg: %s\n", v.Key, v.Value, v.Avg)
        }
    }
}

func convertToSortedSlice(m map[string][]time.Duration) PairList{
  pl := make(PairList, len(m))
  i := 0
  for k, v := range m {
    avg := time.Duration(0)
    for _, dur := range v {
        avg += dur
    }
    val := avg
    avg = avg / time.Duration(len(v))
    pl[i] = Pair{Key: k, Value: val, Avg: avg}
    i++
  }
  sort.Sort(sort.Reverse(pl))
  return pl
}

type Pair struct {
  Key string
  Value time.Duration
  Avg time.Duration
}

type PairList []Pair

func (p PairList) Len() int { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Value < p[j].Value }
func (p PairList) Swap(i, j int){ p[i], p[j] = p[j], p[i] }
