package evaluator

import (
	"monkey/ast"
	"monkey/object"
)

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return newError("Unknown operator: %s%s", operator, right.Type())
	}
}

func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INT_OBJ {
		return newError("Unknown operator: -%s", right.Type())
	}

	value := right.(*object.IntObject).Value
	return &object.IntObject{Value: -value}
}

func evalInfixExpression(operator string, left, right object.Object) object.Object {
	switch {
	case left.Type() == object.INT_OBJ && right.Type() == object.INT_OBJ:
		return evalIntegerInfixExpression(operator, left, right)
	case operator == "==":
		return nativeBoolToObject(left == right)
	case operator == "!=":
		return nativeBoolToObject(left != right)
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s",
			left.Type(), operator, right.Type())
	default:
		return newError("Unknown operator: %s %s %s",
			left.Type(), operator, right.Type())
	}
}

func evalIntegerInfixExpression(operator string, left, right object.Object) object.Object {
	leftValue := left.(*object.IntObject).Value
	rightValue := right.(*object.IntObject).Value

	switch operator {
	case "+":
		return &object.IntObject{Value: leftValue + rightValue}
	case "-":
		return &object.IntObject{Value: leftValue - rightValue}
	case "*":
		return &object.IntObject{Value: leftValue * rightValue}
	case "/":
		return &object.IntObject{Value: leftValue / rightValue}
	case "<":
		return nativeBoolToObject(leftValue < rightValue)
	case ">":
		return nativeBoolToObject(leftValue > rightValue)
	case "==":
		return nativeBoolToObject(leftValue == rightValue)
	case "!=":
		return nativeBoolToObject(leftValue != rightValue)
	default:
		return newError("Unknown operator: %s %s %s",
			left.Type(), operator, right.Type())
	}
}

func evalIfExpression(
	ie *ast.IfExpression,
	env *object.Environment,
) object.Object {
	condition := Eval(ie.Condition, env)
	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return Eval(ie.Consequence, env)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative, env)
	} else {
		return NULL
	}
}

func evalIdentifier(i *ast.Identifier, env *object.Environment) object.Object {
	value, ok := env.Get(i.Value)

	if !ok {
		return newError("Identifier not found: %s", i.Value)
	}

	return value
}

func evalExpressions(exprs []ast.ExpressionNode, env *object.Environment) []object.Object {
	var result []object.Object

	for _, expr := range exprs {
		evaluated := Eval(expr, env)

		if isError(evaluated) {
			return []object.Object{evaluated}
		}

		result = append(result, evaluated)
	}

	return result
}

func applyFunction(fn object.Object, args []object.Object) object.Object {
	fnCasted, ok := fn.(*object.FunctionObject)
	if !ok {
		return newError("Not a function: %s", fn.Type())
	}

	extendedEnv := extendFnEnv(fnCasted, args)
	evaluated := Eval(fnCasted.Body, extendedEnv)

	return unwrapReturnValue(evaluated)
}

func extendFnEnv(fn *object.FunctionObject, args []object.Object) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)

	for id, param := range fn.Parameters {
		env.Set(param.Value, args[id])
	}

	return env
}

func unwrapReturnValue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.ReturnValueObject); ok {
		return returnValue.Value
	}

	return obj
}
