package parser

import (
    "strconv"
    "strings"
    "errors"
    "log"
)

func ParseTokens(tokens [][]Token) ([]INode, error) {
    log.Println("Lexing tokens...")

    nodes := make([]INode, len(tokens))
    
    for ii, line := range tokens {
        var newNode INode

        keyword := line[0]
        var parameter []IParameter

        var lastP IParameter
        for _, p := range line[1:] {
            _, err := strconv.ParseFloat(p.Text, 64)
            if err == nil {
                lastP = NumberParameter{Parameter: Parameter{Text: p.Text}}
                parameter = append(parameter, lastP)
                continue
            }

            //_, isVar := lastP.(VariableParameter)
            _, isLit := lastP.(LiteralParameter)
            //_, isOp := lastP.(OperatorParameter)

            if (isLit && strings.Index(parameter[len(parameter)-1].GetRaw(), "\"") < 1 && parameter[len(parameter)-1].GetRaw() != "\"") || strings.Index(p.Text, "\"") == 0 {
                lastP = LiteralParameter{Parameter: Parameter{Text: p.Text}}
                parameter = append(parameter, lastP)
                continue
            }
            if isOperator(p.Text) {
                lastP = OperatorParameter{Parameter: Parameter{Text: p.Text}}
                parameter = append(parameter, lastP)
            } else {
                lastP = VariableParameter{Parameter: Parameter{Text: p.Text}}
                parameter = append(parameter, lastP)
            }
        }

        if keyword.Text == "end" {
            newNode = &BlockEndNode{Node: Node{Keyword: keyword.Text, Parameter: parameter, ID: ii}}
        } else if keyword.Text == "while" {
            newNode = &LoopNode{Node: Node{Keyword: keyword.Text, Parameter: parameter, ID: ii}}
        } else if keyword.Text == "if" {
            newNode = &ConditionNode{Node: Node{Keyword: keyword.Text, Parameter: parameter, ID: ii}}
        } else if keyword.Text == "number" {
            newNode = &NumberDeclarationNode{Node: Node{Keyword: keyword.Text, Parameter: parameter, ID: ii}}
        } else if keyword.Text == "string" {
            newNode = &LiteralDeclarationNode{Node: Node{Keyword: keyword.Text, Parameter: parameter, ID: ii}}
        } else if keyword.Text == "out" {
            newNode = &OutputNode{Node: Node{Keyword: keyword.Text, Parameter: parameter, ID: ii}}
        } else if keyword.Text == "in" {
            newNode = &InputNode{Node: Node{Keyword: keyword.Text, Parameter: parameter, ID: ii}}
        } else {
            newNode = &SetNode{Node: Node{Keyword: keyword.Text, Parameter: parameter, ID: ii}}
        }

        if ii > 0 {
            calcPrevNext(ii, nodes, newNode)
        }

        if newNode == nil {
            return nil, errors.New("Node type " + keyword.Text + " unknown, maybe a keyword is wrong?")
        }

        nodes[ii] = newNode
    }

    if len(nodes) > 1 {
        lastNode := nodes[len(nodes) - 1]
        if lastNode.GetKeyword() == "end" {
            val := lastNode.(*BlockEndNode)
            val.PreviousNode = nodes[len(nodes) - 2]
        }
    }

    log.Println("Tokens processed, program ready for execution.")

    return nodes, nil
}

func isOperator(t string) bool {
    return t == "==" || t == "=" || t == "!=" ||
        t == "<" || t == ">" || t == "<=" || t == ">=" ||
        t == "+" || t == "-" || t == "/" || t == "*" || t == "%" ||
        t == "(" || t == ")" || t == "**"
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
    default:
        val := nodes[ii - 1].(*SetNode)
        val.NextNode = newNode
        if ii > 1 {
            val.PreviousNode = nodes[ii - 2]
        }
    }
}