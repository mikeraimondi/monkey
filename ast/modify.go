package ast

type ModifierFunc func(Node) Node

func Modify(node Node, modifier ModifierFunc) Node {
	// TODO update 'Token' of parent nodes
	switch node := node.(type) {
	case *Program:
		modifyStatements(node.Statements, modifier)
	case *ExpressionStatement:
		// TODO error handling
		node.Expression, _ = Modify(node.Expression, modifier).(Expression)
	case *InfixExpression:
		// TODO error handling
		node.Left, _ = Modify(node.Left, modifier).(Expression)
		node.Right, _ = Modify(node.Right, modifier).(Expression)
	case *PrefixExpression:
		// TODO error handling
		node.Right, _ = Modify(node.Right, modifier).(Expression)
	case *IndexExpression:
		// TODO error handling
		node.Left, _ = Modify(node.Left, modifier).(Expression)
		node.Index, _ = Modify(node.Index, modifier).(Expression)
	case *IfExpression:
		// TODO error handling
		node.Condition, _ = Modify(node.Condition, modifier).(Expression)
		node.Consequence, _ = Modify(node.Consequence, modifier).(*BlockStatement)
		if node.Alternative != nil {
			node.Alternative, _ = Modify(node.Alternative, modifier).(*BlockStatement)
		}
	case *BlockStatement:
		modifyStatements(node.Statements, modifier)
	case *ReturnStatement:
		// TODO error handling
		node.ReturnValue, _ = Modify(node.ReturnValue, modifier).(Expression)
	case *LetStatement:
		// TODO error handling
		node.Value, _ = Modify(node.Value, modifier).(Expression)
	case *FunctionLiteral:
		// TODO error handling
		for i, _ := range node.Parameters {
			node.Parameters[i], _ = Modify(node.Parameters[i], modifier).(*Identifier)
		}
		node.Body, _ = Modify(node.Body, modifier).(*BlockStatement)
	case *ArrayLiteral:
		// TODO error handling
		for i, _ := range node.Elements {
			node.Elements[i], _ = Modify(node.Elements[i], modifier).(Expression)
		}
	case *HashLiteral:
		// TODO error handling
		newPairs := make(map[Expression]Expression)
		for k, v := range node.Pairs {
			newKey, _ := Modify(k, modifier).(Expression)
			newVal, _ := Modify(v, modifier).(Expression)
			newPairs[newKey] = newVal
		}
		node.Pairs = newPairs
	}

	return modifier(node)
}

func modifyStatements(stmts []Statement, modifier ModifierFunc) {
	// TODO error handling
	for i, statement := range stmts {
		stmts[i], _ = Modify(statement, modifier).(Statement)
	}
}
