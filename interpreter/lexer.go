package interpreter

import (
    "strconv"
    "strings"
    "errors"
    "log"
    "fmt"
)

func ParseTokens(tokens [][]Token) ([]INode, error) {
    log.Println("Lexing tokens...")

    nodes := make([]INode, len(tokens))

    scopeStack := NewScopeStack()
    scopeStack.Push(NewScope(DummyScope))
    
    for ii, line := range tokens {
        var newNode INode

        keyword := line[0]
        var parameter []IParameter

        trace := fmt.Sprintf("File: %s / Line: %d / %s", line[0].File, line[0].Line, keyword.Text)

        var lastP IParameter
        var addedLit bool
        for _, p := range line[1:] {
            _, err := strconv.ParseFloat(p.Text, 64)
            if err == nil {
                lastP = &NumberParameter{Parameter: Parameter{Text: p.Text}}
                parameter = append(parameter, lastP)
                continue
            }

            lp, isLit := lastP.(*LiteralParameter)
            if (!addedLit && isLit && strings.Index(parameter[len(parameter)-1].GetRaw(), "\"") < 1 && parameter[len(parameter)-1].GetRaw() != "\"") ||
                (addedLit && strings.Count(parameter[len(parameter)-1].GetRaw(), "\"") != 2) ||
                strings.Index(p.Text, "\"") == 0 {
                if isLit {
                    lp.Text += " " + p.Text
                    addedLit = true
                } else {
                    lastP = &LiteralParameter{Parameter: Parameter{Text: p.Text}}
                    parameter = append(parameter, lastP)
                    addedLit = false
                }
                continue
            }

            addedLit = false

            if isOperator(p.Text) {
                lastP = &OperatorParameter{Parameter: Parameter{Text: p.Text}}
                parameter = append(parameter, lastP)
            } else {
                lastP = &VariableParameter{Parameter: Parameter{Text: p.Text}}
                parameter = append(parameter, lastP)
            }
        }

        if keyword.Text == "end" {
            newNode = &BlockEndNode{Node: Node{Keyword: keyword.Text, Parameter: parameter, ID: ii, Trace: trace, Scope: scopeStack.Top()}}
            scopeStack.Pop()
        } else if keyword.Text == "while" {
            newNode = &LoopNode{Node: Node{Keyword: keyword.Text, Parameter: parameter, ID: ii, Trace: trace, Scope: scopeStack.Top()}}
            scopeStack.Push(NewScope(scopeStack.Top()))
        } else if keyword.Text == "if" {
            newNode = &ConditionNode{Node: Node{Keyword: keyword.Text, Parameter: parameter, ID: ii, Trace: trace, Scope: scopeStack.Top()}}
            scopeStack.Push(NewScope(scopeStack.Top()))
        } else if keyword.Text == "number" {
            newNode = &NumberDeclarationNode{Node: Node{Keyword: keyword.Text, Parameter: parameter, ID: ii, Trace: trace, Scope: scopeStack.Top()}}
            if len(parameter) != 1 {
                return nil, errors.New(trace + ": Invalid number syntax")
            }
            scopeStack.Top().variableTable[parameter[0].GetRaw()] = float64(0)
        } else if keyword.Text == "string" {
            newNode = &LiteralDeclarationNode{Node: Node{Keyword: keyword.Text, Parameter: parameter, ID: ii, Trace: trace, Scope: scopeStack.Top()}}
            if len(parameter) != 1 {
                return nil, errors.New(trace + ": Invalid string syntax")
            }
            scopeStack.Top().variableTable[parameter[0].GetRaw()] = ""
        } else if keyword.Text == "out" {
            newNode = &OutputNode{Node: Node{Keyword: keyword.Text, Parameter: parameter, ID: ii, Trace: trace, Scope: scopeStack.Top()}}
        } else if keyword.Text == "in" {
            newNode = &InputNode{Node: Node{Keyword: keyword.Text, Parameter: parameter, ID: ii, Trace: trace, Scope: scopeStack.Top()}}
        } else if keyword.Text == "function" {
            if len(parameter) < 1 {
                return nil, errors.New(trace + ": A function declaration needs at least a name as a first parameter")
            }

            passers := make([]Passer, (len(parameter) - 1) / 2)

            counter := 0
            for i := 1; i < len(parameter); i+=2 {
                passers[counter] = Passer{Type: parameter[i].GetRaw(), Name: parameter[i + 1].GetRaw()}
                counter++
            }

            newNode = &FunctionDeclarationNode{Node: Node{Keyword: keyword.Text, Parameter: parameter, ID: ii, Trace: trace, Scope: scopeStack.Top()}, Parameters: passers}

            scopeStack.Top().functionTable[parameter[0].GetRaw()] = newNode

            scopeStack.Push(NewScope(scopeStack.Top()))
        } else if keyword.Text == "call" {
            if len(parameter) < 1 {
                return nil, errors.New(trace + ": A function call needs a function name as a first parameter")
            }

            funcNode := scopeStack.Top().GetFunctionNode(parameter[0].GetRaw())

            if funcNode == nil {
                return nil, errors.New(trace + ": Tried to call invalid function")
            }

            fn := funcNode.(*FunctionDeclarationNode)

            if len(fn.Parameters) != len(parameter) - 1 {
                return nil, errors.New(trace + ": Parameter mismatch")
            }

            passers := make(map[string]IParameter)

            for i := 1; i < len(parameter); i++ {
                passers[fn.Parameters[i - 1].Name] = parameter[i]
            }

            newNode = &CallNode{Node: Node{Keyword: keyword.Text, Parameter: parameter, ID: ii, Trace: trace, Scope: scopeStack.Top()}, Passers: passers}
        } else {
            isVar := scopeStack.Top().GetVar(keyword.Text)

            if isVar != nil {
                newNode = &SetNode{Node: Node{Keyword: keyword.Text, Parameter: parameter, ID: ii, Trace: trace, Scope: scopeStack.Top()}}
            }
        }

        if ii > 0 {
            calcPrevNext(ii, nodes, newNode)
        }

        if newNode == nil {
            return nil, errors.New(trace + ": Node type " + keyword.Text + " unknown, maybe a keyword is wrong? Also check variable declarations/scopes.")
        }

        nodes[ii] = newNode
    }

    // Ugly hack
    if len(nodes) > 1 {
        lastNode := nodes[len(nodes) - 1]
        if lastNode.GetKeyword() == "end" {
            val := lastNode.(*BlockEndNode)
            val.PreviousNode = nodes[len(nodes) - 2]
        }
    }

    log.Println("Initializing nodes...")

    for _, node := range nodes {
        err := node.Init(nodes)
        if err != nil {
            log.Println(node.GetTrace() + ": " + err.Error())
            return nil, errors.New("Error while initializing")
        }
    }

    log.Printf("Tokens processed, %d nodes created. Program ready for execution.\n", len(nodes))

    return nodes, nil
}

func isOperator(t string) bool {
    return t == "==" || t == "=" || t == "!=" ||
        t == "<" || t == ">" || t == "<=" || t == ">=" ||
        t == "+" || t == "-" || t == "/" || t == "*" || t == "%" ||
        t == "(" || t == ")"
}

func calcPrevNext(ii int, nodes []INode, newNode INode) {
    switch nodes[ii - 1].GetKeyword() {
    case "end":
        val := nodes[ii - 1].(*BlockEndNode)
        val.NextNode = newNode
        if ii > 1 {
            val.PreviousNode = nodes[ii - 2]
        }
    case "while":
        val := nodes[ii - 1].(*LoopNode)
        val.NextNode = newNode
        if ii > 1 {
            val.PreviousNode = nodes[ii - 2]
        }
    case "if":
        val := nodes[ii - 1].(*ConditionNode)
        val.NextNode = newNode
        if ii > 1 {
            val.PreviousNode = nodes[ii - 2]
        }
    case "number":
        val := nodes[ii - 1].(*NumberDeclarationNode)
        val.NextNode = newNode
        if ii > 1 {
            val.PreviousNode = nodes[ii - 2]
        }
    case "string":
        val := nodes[ii - 1].(*LiteralDeclarationNode)
        val.NextNode = newNode
        if ii > 1 {
            val.PreviousNode = nodes[ii - 2]
        }
    case "out":
        val := nodes[ii - 1].(*OutputNode)
        val.NextNode = newNode
        if ii > 1 {
            val.PreviousNode = nodes[ii - 2]
        }
    case "in":
        val := nodes[ii - 1].(*InputNode)
        val.NextNode = newNode
        if ii > 1 {
            val.PreviousNode = nodes[ii - 2]
        }
    case "function":
        val := nodes[ii - 1].(*FunctionDeclarationNode)
        val.NextNode = newNode
        if ii > 1 {
            val.PreviousNode = nodes[ii - 2]
        }
    case "call":
        val := nodes[ii - 1].(*CallNode)
        val.NextNode = newNode
        if ii > 1 {
            val.PreviousNode = nodes[ii - 2]
        }
    default:
        val := nodes[ii - 1].(*SetNode)
        val.NextNode = newNode
        if ii > 1 {
            val.PreviousNode = nodes[ii - 2]
        }
    }
}