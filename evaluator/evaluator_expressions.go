package evaluator

import (
	"monkey/ast"
	"monkey/object"
)

func evalPrefixExpression(operator string, rightObject object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(rightObject)
	case "-":
		return evalMinusPrefixOperatorExpression(rightObject)
	default:
		return newError("Unknown operator: %s%s", operator, rightObject.Type())
	}
}

func evalBangOperatorExpression(rightObject object.Object) object.Object {
	switch rightObject {
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

func evalMinusPrefixOperatorExpression(rightObject object.Object) object.Object {
	if rightObject.Type() != object.INT_OBJ {
		return newError("Unknown operator: -%s", rightObject.Type())
	}

	value := rightObject.(*object.IntObject).Value
	return &object.IntObject{Value: -value}
}

func evalInfixExpression(operator string, leftObject, rightObject object.Object) object.Object {
	switch {
	case leftObject.Type() == object.INT_OBJ && rightObject.Type() == object.INT_OBJ:
		return evalIntegerInfixExpression(operator, leftObject, rightObject)
	case operator == "==":
		return nativeBoolToObject(leftObject == rightObject)
	case operator == "!=":
		return nativeBoolToObject(leftObject != rightObject)
	case leftObject.Type() != rightObject.Type():
		return newError("type mismatch: %s %s %s",
			leftObject.Type(), operator, rightObject.Type())
	default:
		return newError("Unknown operator: %s %s %s",
			leftObject.Type(), operator, rightObject.Type())
	}
}

func evalIntegerInfixExpression(operator string, leftObject, rightObject object.Object) object.Object {
	leftValue := leftObject.(*object.IntObject).Value
	rightValue := rightObject.(*object.IntObject).Value

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
			leftObject.Type(), operator, rightObject.Type())
	}
}

func evalIfExpression(
	node *ast.IfExpressionNode,
	env *object.Environment,
) object.Object {
	condition := Eval(node.Condition, env)
	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return Eval(node.Consequence, env)
	} else if node.Alternative != nil {
		return Eval(node.Alternative, env)
	} else {
		return NULL
	}
}

func evalIdentifier(node *ast.IdentifierNode, env *object.Environment) object.Object {
	value, ok := env.Get(node.Value)

	if !ok {
		return newError("Identifier not found: %s", node.Value)
	}

	return value
}

func evalExpressions(exprNodes []ast.ExpressionNode, env *object.Environment) []object.Object {
	var result []object.Object

	for _, expr := range exprNodes {
		evaluated := Eval(expr, env)

		if isError(evaluated) {
			return []object.Object{evaluated}
		}

		result = append(result, evaluated)
	}

	return result
}

func applyFunction(fnObject object.Object, args []object.Object) object.Object {
	fnCasted, ok := fnObject.(*object.FunctionObject)
	if !ok {
		return newError("Not a function: %s", fnObject.Type())
	}

	extendedEnv := extendFnEnv(fnCasted, args)
	evaluated := Eval(fnCasted.Body, extendedEnv)

	return unwrapReturnValue(evaluated)
}

func extendFnEnv(fnObject *object.FunctionObject, args []object.Object) *object.Environment {
	env := object.NewEnclosedEnvironment(fnObject.Env)

	for id, param := range fnObject.Parameters {
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
