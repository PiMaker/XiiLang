package interpreter


func NewScopeStack() *ScopeStack {
	return &ScopeStack{}
}

type ScopeStack struct {
	nodes []*Scope
	count int
}

func (s *ScopeStack) Push(n *Scope) {
	s.nodes = append(s.nodes[:s.count], n)
	s.count++
}

func (s *ScopeStack) Pop() *Scope {
	s.count--
	return s.nodes[s.count]
}

func (s *ScopeStack) Top() *Scope {
    return s.nodes[s.count - 1]
}


func NewNodeStack() *NodeStack {
	return &NodeStack{}
}

type NodeStack struct {
	nodes []INode
	count int
}

func (s *NodeStack) Push(n INode) {
	s.nodes = append(s.nodes[:s.count], n)
	s.count++
}

func (s *NodeStack) Pop() INode {
	s.count--
	return s.nodes[s.count]
}

func (s *NodeStack) Top() INode {
    return s.nodes[s.count - 1]
}