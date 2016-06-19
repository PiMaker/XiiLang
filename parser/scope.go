package parser

import (
)

type Scope struct {
    baseScope *Scope
    variableTable map[string]interface{}
    functionTable map[string]INode
}

var DummyScope = &Scope{}

func NewScope(baseScope *Scope) *Scope {
    return &Scope{variableTable: make(map[string]interface{}), functionTable: make(map[string]INode), baseScope: baseScope}
}

func (scope *Scope) SetVar(name string, value interface{}) {
    if !scope.setIfExists(name, value) {
        scope.variableTable[name] = value
    }
}

func (scope *Scope) setIfExists(name string, value interface{}) bool {
    _, ok := scope.variableTable[name]
    if ok {
        scope.variableTable[name] = value
        return true
    }

    if scope.baseScope != nil {
        return scope.baseScope.setIfExists(name, value)
    }

    return false
}

func (scope *Scope) GetVar(name string) interface{} {
    val, ok := scope.variableTable[name]
    if ok {
        return val
    }

    if scope.baseScope != nil {
        return scope.baseScope.GetVar(name)
    }

    return nil
}

func (scope *Scope) GetFunctionNode(name string) INode {
    val, ok := scope.functionTable[name]
    if ok {
        return val
    }

    if scope.baseScope != nil {
        return scope.baseScope.GetFunctionNode(name)
    }

    return nil
}