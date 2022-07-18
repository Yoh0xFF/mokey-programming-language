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
		return Eval(node.Expression, env)

	case *ast.ReturnStatementNode:
		value := Eval(node.ReturnValue, env)
		if isError(value) {
			return value
		}
		return &object.ReturnValueObject{Value: value}

	case *ast.LetStatementNode:
		value := Eval(node.Value, env)
		if isError(value) {
			return value
		}
		env.Set(node.Name.Value, value)

	// Expressions
	case *ast.IntegerLiteral:
		return &object.IntObject{Value: node.Value}

	case *ast.Boolean:
		return nativeBoolToObject(node.Value)

	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)

	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}

		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}

		return evalInfixExpression(node.Operator, left, right)

	case *ast.IfExpression:
		return evalIfExpression(node, env)

	case *ast.Identifier:
		return evalIdentifier(node, env)

	case *ast.FunctionLiteral:
		return &object.FunctionObject{
			Parameters: node.Parameters,
			Env:        env,
			Body:       node.Body,
		}

	case *ast.CallExpression:
		fn := Eval(node.Function, env)
		if isError(fn) {
			return fn
		}

		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}

		return applyFunction(fn, args)
	}

	return nil
}
