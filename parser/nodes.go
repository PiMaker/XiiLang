package parser

import (
    "strconv"
    "errors"
    "fmt"
)

type Node struct {
    Keyword string
    ID int
    Parameter []IParameter
    NextNode, PreviousNode INode
}

type INode interface {
    Execute(state *XiiState) error
    Previous() INode
    Next() INode
    GetKeyword() string
    GetID() int
}


type NumberDeclarationNode struct {
    Node
}

func (node NumberDeclarationNode) Execute(state *XiiState) error {
    if len(node.Parameter) != 1 {
        return errors.New("number: No parameter name given (or too many)")
    }

    state.VariableNumberTable[node.Parameter[0].GetRaw()] = 0

    return nil
}


type LiteralDeclarationNode struct {
    Node
}

func (node LiteralDeclarationNode) Execute(state *XiiState) error {
    if len(node.Parameter) != 1 {
        return errors.New("string: No parameter name given (or too many)")
    }

    state.VariableLiteralTable[node.Parameter[0].GetRaw()] = ""

    return nil
}


type OutputNode struct {
    Node
}

func (node OutputNode) Execute(state *XiiState) error {
    for i, n := range node.Parameter {
        _, ok := n.(VariableParameter)
        if i != 0 && !ok {
            fmt.Print(" ")
        }
        fmt.Print(n.GetText(*state))
    }

    fmt.Println()

    return nil
}


type InputNode struct {
    Node
}

func (node InputNode) Execute(state *XiiState) error {
    if len(node.Parameter) != 1 {
        return errors.New("in: No parameter name given (or too many)")
    }

    varname := node.Parameter[0].GetRaw()
    _, ok1 := state.VariableLiteralTable[varname]
    _, ok2 := state.VariableNumberTable[varname]

    if ok1 || ok2 {
        var text string
        fmt.Scanln(&text)
        if ok1 {
            state.VariableLiteralTable[varname] = text
        } else if ok2 {
            for {
                num, err := strconv.ParseFloat(text, 64)
                if err == nil {
                    state.VariableNumberTable[varname] = num
                    break
                } else {
                    fmt.Println("Please retry: " + err.Error())
                    fmt.Scanln(&text)
                }
            }
        }
    } else {
        return errors.New("in: Unknown variable cannot be read into")
    }

    return nil
}


type LoopNode struct {
    Node
}

func (node LoopNode) Execute(state *XiiState) error {
    if Evaluate(*state, node.Parameter) == 0 {
        nextEnd := node.Next()
        _, ok := nextEnd.(*BlockEndNode)
        for !ok {
            nextEnd = nextEnd.Next()
            _, ok = nextEnd.(*BlockEndNode)
        }
        state.NextNode = nextEnd.Next()
    }

    return nil
}


type ConditionNode struct {
    Node
}

func (node ConditionNode) Execute(state *XiiState) error {
    if Evaluate(*state, node.Parameter) == 0 {
        nextEnd := node.Next()
        _, ok := nextEnd.(*BlockEndNode)
        for !ok {
            nextEnd = nextEnd.Next()
            _, ok = nextEnd.(*BlockEndNode)
        }
        state.NextNode = nextEnd.Next()
    }

    return nil
}


type BlockEndNode struct {
    Node
}

func (node BlockEndNode) Execute(state *XiiState) error {
    companion := node.Previous()
    for {
        switch companion.(type) {
        case (*ConditionNode):
            return nil
        case (*LoopNode):
            state.NextNode = companion
            return nil
        }
        companion = companion.Previous()
    }
}


type SetNode struct {
    Node
}

func (node SetNode) Execute(state *XiiState) error {
    if len(node.Parameter) < 2 || node.Parameter[0].GetRaw() != "=" {
        return errors.New("set: Invalid set syntax")
    }

    varname := node.Keyword
    _, ok := state.VariableNumberTable[varname]

    if !ok {
        _, ok = state.VariableLiteralTable[varname]
        if !ok {
            return errors.New("set: Can't set not existing variable")
        }

        var setval string
        for i, v := range node.Parameter {
            if i > 0 {
                setval += " "
            }
            setval += v.GetText(*state)
        }

        state.VariableLiteralTable[varname] = setval

        return nil
    }

    state.VariableNumberTable[varname] = Evaluate(*state, node.Parameter[1:])

    return nil
}


func (node Node) Execute(state *XiiState) error {
    return errors.New("Error: NoOp Node executed")
}

func (node Node) Previous() INode {
    return node.PreviousNode
}

func (node Node) Next() INode {
    return node.NextNode
}

func (node Node) GetKeyword() string {
    return node.Keyword
}

func (node Node) GetID() int {
    return node.ID
}

func (node Node) String() string {
    return fmt.Sprintf("{{%d/%s : %s}}", node.ID, node.Keyword, node.Parameter)
}