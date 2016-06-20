package interpreter

import (
	"strings"
	"errors"
	"fmt"

	"github.com/Knetic/govaluate"
)

var VerboseEval bool

type Expression struct {
	Expr *govaluate.EvaluableExpression
	ExprString string
}

func Evaluate(node INode, expression *Expression) (float64, error) {

	result, err := expression.Expr.Evaluate(node.GetScope())

	if VerboseEval {
		fmt.Printf("Evaluated expression: %s -> %s\n", expression.ExprString, result)
	}

	if err != nil {
		return 0, err
	}

	switch v := result.(type) {
	case float64:
		return v, nil
	case bool:
		if v {
			return float64(1), nil
		}
		return float64(0), nil
	}

	return 0, errors.New("Unexpected expression evaluation result")
}

func NewExpression(condition []IParameter) (*Expression, error) {
	var str string
	for _, param := range condition {
		str += param.GetRaw() + " "
	}

	if strings.TrimSpace(str) == "" {
		return nil, errors.New("Empty expression passed")
	}

	expr, err := govaluate.NewEvaluableExpression(str)

	if err != nil {
		return nil, err
	}

	if VerboseEval {
		fmt.Println("Created expression: " + str)
	}

	return &Expression{Expr: expr, ExprString: str}, nil
}