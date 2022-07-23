package evaluator

import (
	"monkey/ast"
	"monkey/object"
)

var (
	NULL  = &object.NullObject{}
	TRUE  = &object.BoolObject{Value: true}
	FALSE = &object.BoolObject{Value: false}
)

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	// Statements
	case *ast.ProgramNode:
		return evalProgram(node, env)

	case *ast.BlockStatementNode:
		return evalBlockStatement(node, env)

	case *ast.ExpressionStatementNode:
		return Eval(node.ExpressionNode, env)

	case *ast.ReturnStatementNode:
		resultObject := Eval(node.ReturnValueNode, env)
		if isError(resultObject) {
			return resultObject
		}
		return &object.ReturnValueObject{ValueObject: resultObject}

	case *ast.LetStatementNode:
		resultObject := Eval(node.ValueNode, env)
		if isError(resultObject) {
			return resultObject
		}
		env.Set(node.NameNode.Value, resultObject)

	// Expressions
	case *ast.IntegerLiteralNode:
		return &object.IntObject{Value: node.Value}

	case *ast.StringLiteralNode:
		return &object.StringObject{Value: node.Value}

	case *ast.BooleanNode:
		return nativeBoolToObject(node.Value)

	case *ast.PrefixExpressionNode:
		rightObject := Eval(node.RightNode, env)
		if isError(rightObject) {
			return rightObject
		}
		return evalPrefixExpression(node.Operator, rightObject)

	case *ast.InfixExpressionNode:
		leftObject := Eval(node.LeftNode, env)
		if isError(leftObject) {
			return leftObject
		}

		rightObject := Eval(node.RightNode, env)
		if isError(rightObject) {
			return rightObject
		}

		return evalInfixExpression(node.Operator, leftObject, rightObject)

	case *ast.IfExpressionNode:
		return evalIfExpression(node, env)

	case *ast.IdentifierNode:
		return evalIdentifier(node, env)

	case *ast.FunctionLiteralNode:
		return &object.FunctionObject{
			ParamNodes: node.ParamNodes,
			Env:        env,
			BodyNode:   node.BodyNode,
		}

	case *ast.CallExpressionNode:
		fnObject := Eval(node.FnNode, env)
		if isError(fnObject) {
			return fnObject
		}

		argObjects := evalExpressions(node.ArgNodes, env)
		if len(argObjects) == 1 && isError(argObjects[0]) {
			return argObjects[0]
		}

		return applyFunction(fnObject, argObjects)

	case *ast.ArrayLiteralNode:
		elementObjects := evalExpressions(node.Elements, env)

		if len(elementObjects) == 1 && isError(elementObjects[0]) {
			return elementObjects[0]
		}

		return &object.ArrayObject{Elements: elementObjects}

	case *ast.IndexExpressionNode:
		leftObject := Eval(node.Left, env)
		if isError(leftObject) {
			return leftObject
		}

		indexObject := Eval(node.Index, env)
		if isError(indexObject) {
			return indexObject
		}

		return evalIndexExpression(leftObject, indexObject)

	case *ast.HashLiteralNode:
		return evalHashLiteral(node, env)
	}

	return nil
}
