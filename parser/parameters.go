package parser

import (
	"strings"

    humanize "github.com/dustin/go-humanize"
)

type Parameter struct {
    Text string
}

type IParameter interface {
    GetText(state XiiState) string
    GetRaw() string
}

type LiteralParameter struct {
    Parameter
}

func (p LiteralParameter) String() string {
    return "'" + p.Text + "'"
}

type NumberParameter struct {
    Parameter
}

func (p NumberParameter) String() string {
    return "*" + p.Text + "*"
}

type VariableParameter struct {
    Parameter
}

func (p VariableParameter) String() string {
    return "%" + p.Text + "%"
}

func (p VariableParameter) GetText(state XiiState) string {
    lit, ok1 := state.VariableTable[p.Text].(string)
    num, ok2 := state.VariableTable[p.Text].(float64)

    if ok1 {
        return lit
    }
    if ok2 {
        return humanize.Ftoa(num)
    }

    return ""
}

func (l LiteralParameter) GetText(state XiiState) string {
    //if l.Text == "\"" {
    //    return " "
    //}
    return strings.Replace(l.Text, "\"", "", -1)
}

type OperatorParameter struct {
    Parameter
}


func (p OperatorParameter) String() string {
    return "{" + p.Text + "}"
}

func (p Parameter) GetText(state XiiState) string {
    return p.Text
}

func (p Parameter) GetRaw() string {
    return p.Text
}