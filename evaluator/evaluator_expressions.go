package evaluator

import "monkey/object"

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

	value := right.(*object.Int).Value
	return &object.Int{Value: -value}
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
		return newError("unknown operator: %s %s %s",
			left.Type(), operator, right.Type())
	}
}

func evalIntegerInfixExpression(operator string, left, right object.Object) object.Object {
	leftValue := left.(*object.Int).Value
	rightValue := right.(*object.Int).Value

	switch operator {
	case "+":
		return &object.Int{Value: leftValue + rightValue}
	case "-":
		return &object.Int{Value: leftValue - rightValue}
	case "*":
		return &object.Int{Value: leftValue * rightValue}
	case "/":
		return &object.Int{Value: leftValue / rightValue}
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
