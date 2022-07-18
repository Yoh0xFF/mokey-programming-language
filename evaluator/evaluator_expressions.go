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
		return newErrorObject("Unknown operator: %s%s", operator, rightObject.Type())
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
		return newErrorObject("Unknown operator: -%s", rightObject.Type())
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
		return newErrorObject("type mismatch: %s %s %s",
			leftObject.Type(), operator, rightObject.Type())
	default:
		return newErrorObject("Unknown operator: %s %s %s",
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
		return newErrorObject("Unknown operator: %s %s %s",
			leftObject.Type(), operator, rightObject.Type())
	}
}

func evalIfExpression(
	node *ast.IfExpressionNode,
	env *object.Environment,
) object.Object {
	conditionObject := Eval(node.ConditionNode, env)
	if isError(conditionObject) {
		return conditionObject
	}

	if isTruthy(conditionObject) {
		return Eval(node.ConsequenceNode, env)
	} else if node.AlternativeNode != nil {
		return Eval(node.AlternativeNode, env)
	} else {
		return NULL
	}
}

func evalIdentifier(node *ast.IdentifierNode, env *object.Environment) object.Object {
	value, ok := env.Get(node.Value)

	if !ok {
		return newErrorObject("Identifier not found: %s", node.Value)
	}

	return value
}

func evalExpressions(exprNodes []ast.ExpressionNode, env *object.Environment) []object.Object {
	var resultObjects []object.Object

	for _, exprNode := range exprNodes {
		resultObject := Eval(exprNode, env)

		if isError(resultObject) {
			return []object.Object{resultObject}
		}

		resultObjects = append(resultObjects, resultObject)
	}

	return resultObjects
}

func applyFunction(fnObject object.Object, argObjects []object.Object) object.Object {
	fnObjectCasted, ok := fnObject.(*object.FunctionObject)
	if !ok {
		return newErrorObject("Not a function: %s", fnObject.Type())
	}

	extendedEnv := extendFnEnv(fnObjectCasted, argObjects)
	resultObject := Eval(fnObjectCasted.BodyNode, extendedEnv)

	return unwrapReturnValue(resultObject)
}

func extendFnEnv(fnObject *object.FunctionObject, argObjects []object.Object) *object.Environment {
	env := object.NewEnclosedEnvironment(fnObject.Env)

	for id, paramNode := range fnObject.ParamNodes {
		env.Set(paramNode.Value, argObjects[id])
	}

	return env
}

func unwrapReturnValue(obj object.Object) object.Object {
	if returnValueObject, ok := obj.(*object.ReturnValueObject); ok {
		return returnValueObject.ValueObject
	}

	return obj
}
