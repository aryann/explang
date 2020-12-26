package evaluator

import (
	"explang/parser"
	"fmt"
)

type evaluator struct {
	*environment
}

func newEvaluator() *evaluator {
	return &evaluator{
		environment: newEnvironment(),
	}
}

func (e *evaluator) VisitFile(node parser.File) (parser.Node, error) {
	newFile := &parser.File{Nodes: make([]parser.Node, 0)}
	for _, node := range node.Nodes {
		newNode, err := node.Accept(e)
		if err != nil {
			return nil, err
		}
		newFile.Nodes = append(newFile.Nodes, newNode)
	}
	return newFile, nil
}

func (e *evaluator) VisitDefinition(node parser.Definition) (parser.Node, error) {
	value, err := node.Expression.Accept(e)
	if err != nil {
		return nil, fmt.Errorf("could not evaluate expression: %w", err)
	}
	e.Set(node.Name, value)
	return node, nil
}

func (e *evaluator) VisitString(node parser.String) (parser.Node, error) {
	return node, nil
}

func (e *evaluator) VisitRef(node parser.Ref) (parser.Node, error) {
	val, ok := e.Get(node.Name)
	if !ok {
		return nil, fmt.Errorf("no such variable: %s", node.Name)
	}
	return val, nil
}

func Evaluate(node parser.Node) (parser.Node, error) {
	evaluator := newEvaluator()
	return node.Accept(evaluator)
}