package interpreter

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
    Trace string
    Scope *Scope
}

type INode interface {
    Execute(state *XiiState) error
    Previous() INode
    Next() INode
    Init(nodes []INode) error
    GetKeyword() string
    GetID() int
    GetTrace() string
    GetScope() *Scope
}


type NumberDeclarationNode struct {
    Node
}

func (node *NumberDeclarationNode) Execute(state *XiiState) error {
    return nil
}


type LiteralDeclarationNode struct {
    Node
}

func (node *LiteralDeclarationNode) Execute(state *XiiState) error {
    return nil
}


type Passer struct {
    Name string
    Type string
}


type FunctionDeclarationNode struct {
    Node
    Parameters []Passer
    nextAfterEnd INode
}

func (node *FunctionDeclarationNode) Init(nodes []INode) error {
    nextEnd := findNextEndNode(node)

    if nextEnd == nil {
        return errors.New("A function node requires a matching end node")
    }

    node.nextAfterEnd = nextEnd.Next()

    return nil
}

func (node *FunctionDeclarationNode) Execute(state *XiiState) error {
    if state.PassingArea == nil {
        state.NextNode = node.nextAfterEnd
    } else {
        for k, v := range state.PassingArea {
            node.GetScope().SetVar(k, v)
        }
        state.PassingArea = nil
    }
    
    return nil
}


type CallNode struct {
    Node
    Passers map[string]IParameter
}

func (node *CallNode) Execute(state *XiiState) error {
    fun := node.GetScope().GetFunctionNode(node.Parameter[0].GetRaw())
    if fun == nil {
        return errors.New("Tried to call non-existing function")
    }

    state.PassingArea = make(map[string]interface{}, 0)
    for k, v := range node.Passers {
        state.PassingArea[k] = v.GetValue(node.GetScope())
    }

    state.NextNode = fun
    state.FunctionStack.Push(node)

    return nil
}


type OutputNode struct {
    Node
}

func (node *OutputNode) Execute(state *XiiState) error {
    for i, n := range node.Parameter {
        _, ok := n.(VariableParameter)
        if i != 0 && !ok {
            state.StdOut.WriteRune(' ')
        }
        state.StdOut.WriteString(n.GetText(node.GetScope()))
    }

    state.StdOut.WriteRune('\n')

    state.StdOut.Flush()

    return nil
}


type InputNode struct {
    Node
}

func (node *InputNode) Execute(state *XiiState) error {
    if len(node.Parameter) != 1 {
        return errors.New("in: No parameter name given (or too many)")
    }

    varname := node.Parameter[0].GetRaw()
    variable := node.GetScope().GetVar(varname)

    if variable == nil {
        return errors.New("Tried to 'in' not existing variable")
    }

    _, ok1 := variable.(string)
    _, ok2 := variable.(float64)

    if ok1 || ok2 {
        var text string
        fmt.Scanln(&text)
        if ok1 {
            node.GetScope().SetVar(varname, text)
        } else if ok2 {
            for {
                num, err := strconv.ParseFloat(text, 64)
                if err == nil {
                    node.GetScope().SetVar(varname, num)
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
    nextAfterEndNode INode
    expression *Expression
}

func (node *LoopNode) Init(nodes []INode) error {
    nextEnd := findNextEndNode(node)

    if nextEnd == nil {
        return errors.New("A loop node requires a matching end node")
    }

    node.nextAfterEndNode = nextEnd.Next()

    exp, err := NewExpression(node.Parameter)

    if err != nil {
        return err
    }

    node.expression = exp

    return nil
}

func (node *LoopNode) Execute(state *XiiState) error {
    res, err := Evaluate(node, node.expression)

    if err != nil {
        return err
    }

    if res == 0 {
        state.NextNode = node.nextAfterEndNode
    }

    return nil
}


type ConditionNode struct {
    Node
    nextEndNode INode
    expression *Expression
}

func (node *ConditionNode) Init(nodes []INode) error {
    nextEnd := findNextEndNode(node)

    if nextEnd == nil {
        return errors.New("A condition node requires a matching end node")
    }
    
    node.nextEndNode = nextEnd.Next()

    exp, err := NewExpression(node.Parameter)

    if err != nil {
        return err
    }

    node.expression = exp

    return nil
}

func (node *ConditionNode) Execute(state *XiiState) error {
    res, err := Evaluate(node, node.expression)

    if err != nil {
        return err
    }

    if res == 0 {
        state.NextNode = node.nextEndNode
    }

    return nil
}


type BlockEndNode struct {
    Node
    companionNode INode
    endsFunction bool
}

func (node *BlockEndNode) Init(nodes []INode) error {
    companion := node.Previous()
    counter := 1
    for {
        switch companion.(type) {
        case (*ConditionNode):
            counter--
            if counter == 0 {
                return nil
            }
        case (*LoopNode):
            counter--
            if counter == 0 {
                node.companionNode = companion
                return nil
            }
        case (*FunctionDeclarationNode):
            counter--
            if counter == 0 {
                node.endsFunction = true
                return nil
            }
        case (*BlockEndNode):
            counter++
        }
        companion = companion.Previous()
        if companion == nil {
            return errors.New("end Node without condition/loop")
        }
    }
}

func (node *BlockEndNode) Execute(state *XiiState) error {
    if node.companionNode != nil {
        state.NextNode = node.companionNode
    }

    if node.endsFunction {
        state.NextNode = state.FunctionStack.Pop().Next()
    }

    return nil
}


type SetNode struct {
    Node
    expression *Expression
}

func (node *SetNode) Init(nodes []INode) error {
    if len(node.Parameter) < 2 || node.Parameter[0].GetRaw() != "=" {
        return errors.New("set: Invalid set syntax")
    }

    exp, err := NewExpression(node.Parameter[1:])

    if err != nil {
        return err
    }

    node.expression = exp

    return nil
}

func (node *SetNode) Execute(state *XiiState) error {
    varname := node.Keyword
    variable := node.GetScope().GetVar(varname)
    _, ok := variable.(float64)

    if !ok {
        _, ok = variable.(string)
        if !ok {
            return errors.New("set: Can't set not existing variable")
        }

        var setval string
        for i, v := range node.Parameter {
            if i > 0 {
                setval += " "
            }
            setval += v.GetText(node.GetScope())
        }

        node.GetScope().SetVar(varname, setval)

        return nil
    }
    
    res, err := Evaluate(node, node.expression)

    if err != nil {
        return err
    }

    node.GetScope().SetVar(varname, res)

    return nil
}


func (node *Node) Execute(state *XiiState) error {
    return errors.New("No-Op Node executed")
}

func (node *Node) Previous() INode {
    return node.PreviousNode
}

func (node *Node) Next() INode {
    return node.NextNode
}

func (node *Node) GetKeyword() string {
    return node.Keyword
}

func (node *Node) GetID() int {
    return node.ID
}

func (node *Node) Init(nodes []INode) error {
    return nil
}

func (node *Node) String() string {
    return fmt.Sprintf("{{%d/%s : %s}}", node.ID, node.Keyword, node.Parameter)
}

func (node *Node) GetTrace() string {
    return node.Trace
}

func (node *Node) GetScope() *Scope {
    return node.Scope
}

func findNextEndNode(node INode) INode {
    nextEnd := node.Next()
    counter := 1
    _, ok := nextEnd.(*BlockEndNode)
    for {
        if ok {
            counter--

            if counter == 0 {
                break
            }
        }

        _, ok2 := nextEnd.(*LoopNode)
        _, ok3 := nextEnd.(*ConditionNode)
        _, ok4 := nextEnd.(*FunctionDeclarationNode)

        if ok2 || ok3 || ok4 {
            counter++
        }
        
        nextEnd = nextEnd.Next()

        if nextEnd == nil {
            return nil
        }

        _, ok = nextEnd.(*BlockEndNode)
    }

    return nextEnd
}
