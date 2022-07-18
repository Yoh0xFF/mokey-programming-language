package evaluator

import (
	"monkey/ast"
	"monkey/object"
)

func evalProgram(node *ast.ProgramNode, env *object.Environment) object.Object {
	var resultObject object.Object

	for _, statementNode := range node.StatementNodes {
		resultObject = Eval(statementNode, env)

		switch resultObject := resultObject.(type) {
		case *object.ReturnValueObject:
			return resultObject.ValueObject
		case *object.ErrorObject:
			return resultObject
		}
	}

	return resultObject
}

func evalBlockStatement(node *ast.BlockStatementNode, env *object.Environment) object.Object {
	var resultObject object.Object

	for _, statementNode := range node.StatementNodes {
		resultObject = Eval(statementNode, env)

		if resultObject != nil && (resultObject.Type() == object.RETURN_VALUE_OBJ || resultObject.Type() == object.ERROR_OBJ) {
			return resultObject
		}
	}

	return resultObject
}
