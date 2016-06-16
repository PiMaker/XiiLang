package parser

type XiiState struct {
    VariableTable map[string]interface{}
    FunctionTable map[string]*Node
    NextNode INode
    Nodes []INode
}