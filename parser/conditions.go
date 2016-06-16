package parser

import (
    "github.com/soniah/evaler"
    "fmt"
)

var VerboseEval bool

func Evaluate(state XiiState, condition []IParameter) float64 {
    var toEval string
    var rawEval string
    for _, v := range condition {
        toEval += v.GetText(state)
        rawEval += v.GetRaw()
    }

    if VerboseEval {
        fmt.Printf("Executing eval: %s (%s)\n", rawEval, toEval)
    }

    result, err := evaler.Eval(toEval)

    if err != nil {
        fmt.Println("Expression eval error:")
        fmt.Println(err.Error())
        return 0
    }

    val, _ := result.Float64()

    if VerboseEval {
        fmt.Printf("Eval-Result: %f\n", val)
    }

    return val
}