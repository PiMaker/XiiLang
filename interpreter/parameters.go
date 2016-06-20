package interpreter

import (
	"strings"
    "strconv"

    humanize "github.com/dustin/go-humanize"
)

type Parameter struct {
    Text string
}

type IParameter interface {
    GetText(scope *Scope) string
    GetRaw() string
    GetValue(scope *Scope) interface{}
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

func (p NumberParameter) GetValue(_ *Scope) interface{} {
    retval, _ := strconv.ParseFloat(p.Text, 64)
    return retval
}

type VariableParameter struct {
    Parameter
}

func (p VariableParameter) String() string {
    return "%" + p.Text + "%"
}

func (p VariableParameter) GetText(scope *Scope) string {
    variable := scope.GetVar(p.Text)

    if variable == nil {
        return ""
    }

    lit, ok1 := variable.(string)
    num, ok2 := variable.(float64)

    if ok1 {
        return strings.Replace(lit, "\"", "", -1)
    }
    if ok2 {
        return humanize.Ftoa(num)
    }

    return ""
}

func (p VariableParameter) GetValue(scope *Scope) interface{} {
    variable := scope.GetVar(p.Text)

    if variable == nil {
        return ""
    }

    lit, ok1 := variable.(string)
    num, ok2 := variable.(float64)

    if ok1 {
        return strings.Replace(lit, "\"", "", -1)
    }
    if ok2 {
        return num
    }

    return ""
}

func (l LiteralParameter) GetText(scope *Scope) string {
    return strings.Replace(l.Text, "\"", "", -1)
}

type OperatorParameter struct {
    Parameter
}


func (p OperatorParameter) String() string {
    return "{" + p.Text + "}"
}

func (p Parameter) GetText(_ *Scope) string {
    return p.Text
}

func (p Parameter) GetRaw() string {
    return p.Text
}

func (p Parameter) GetValue(scope *Scope) interface{} {
    return p.GetText(scope)
}