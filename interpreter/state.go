package interpreter

import (
    "bufio"
)

type XiiState struct {
    NextNode INode
    Nodes []INode
    FunctionStack *NodeStack
    PassingArea map[string]interface{}
    StdOut *bufio.Writer
}