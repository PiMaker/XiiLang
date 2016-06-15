package parser

import (
    "github.com/soniah/evaler/stack"
)

type XiiState struct {
    VariableLiteralTable map[string]string
    VariableNumberTable map[string]float64
    FunctionTable map[string]*Node
    FunctionStack stack.Stack
    NextNode INode
    Nodes []INode
}