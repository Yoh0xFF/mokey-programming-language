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
	case leftObject.Type() == object.STRING_OBJ && rightObject.Type() == object.STRING_OBJ:
		return evalStringInfixExpression(operator, leftObject, rightObject)
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

func evalStringInfixExpression(
	operator string,
	left, right object.Object,
) object.Object {
	if operator != "+" {
		return newErrorObject("Unknown operator: %s %s %s",
			left.Type(), operator, right.Type())
	}

	leftVal := left.(*object.StringObject).Value
	rightVal := right.(*object.StringObject).Value
	return &object.StringObject{Value: leftVal + rightVal}
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
	if val, ok := env.Get(node.Value); ok {
		return val
	}

	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}

	return newErrorObject("Identifier not found: " + node.Value)
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
	switch fnObjectCasted := fnObject.(type) {
	case *object.FunctionObject:

		extendedEnv := extendFnEnv(fnObjectCasted, argObjects)
		resultObject := Eval(fnObjectCasted.BodyNode, extendedEnv)
		return unwrapReturnValue(resultObject)

	case *object.BuiltinObject:
		return fnObjectCasted.Fn(argObjects...)

	default:
		return newErrorObject("Not a function: %s", fnObject.Type())
	}
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

func evalIndexExpression(leftObject, indexObject object.Object) object.Object {
	switch {
	case leftObject.Type() == object.ARRAY_OBJ && indexObject.Type() == object.INT_OBJ:
		return evalArrayIndexExpression(leftObject, indexObject)
	case leftObject.Type() == object.HASH_OBJ:
		return evalHashIndexExpression(leftObject, indexObject)
	default:
		return newErrorObject("Index operator not supported: %s", leftObject.Type())
	}
}

func evalArrayIndexExpression(arrayObject, indexObject object.Object) object.Object {
	arrayObjectCasted := arrayObject.(*object.ArrayObject)
	idx := indexObject.(*object.IntObject).Value
	max := int64(len(arrayObjectCasted.Elements) - 1)

	if idx < 0 || idx > max {
		return NULL
	}

	return arrayObjectCasted.Elements[idx]
}

func evalHashLiteral(
	node *ast.HashLiteralNode,
	env *object.Environment,
) object.Object {
	pairs := make(map[object.HashKey]object.HashPair)

	for keyNode, valueNode := range node.Pairs {
		keyObject := Eval(keyNode, env)
		if isError(keyObject) {
			return keyObject
		}

		hashKey, ok := keyObject.(object.Hashable)
		if !ok {
			return newErrorObject("Unusable as hash key: %s", keyObject.Type())
		}

		valueObject := Eval(valueNode, env)
		if isError(valueObject) {
			return valueObject
		}

		hashed := hashKey.HashKey()
		pairs[hashed] = object.HashPair{Key: keyObject, Value: valueObject}
	}

	return &object.HashObject{Pairs: pairs}
}

func evalHashIndexExpression(hashObject, indexObject object.Object) object.Object {
	hashObjectCasted := hashObject.(*object.HashObject)

	key, ok := indexObject.(object.Hashable)
	if !ok {
		return newErrorObject("Unusable as hash key: %s", indexObject.Type())
	}

	pair, ok := hashObjectCasted.Pairs[key.HashKey()]
	if !ok {
		return NULL
	}

	return pair.Value
}
