package evaluator

import (
	"fmt"
	"monkey/object"
)

func nativeBoolToObject(input bool) *object.BoolObject {
	if input {
		return TRUE
	}

	return FALSE
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case NULL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		return true
	}
}

func newErrorObject(format string, a ...interface{}) *object.ErrorObject {
	return &object.ErrorObject{Message: fmt.Sprintf(format, a...)}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}
