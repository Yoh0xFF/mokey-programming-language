package compiler

import (
	"fmt"
	"monkey/ast"
	"monkey/code"
	"monkey/object"
	"sort"
)

func (c *Compiler) Compile(node ast.Node) error {
	switch node := node.(type) {
	case *ast.ProgramNode:
		for _, s := range node.StatementNodes {
			err := c.Compile(s)
			if err != nil {
				return err
			}
		}

	case *ast.ExpressionStatementNode:
		err := c.Compile(node.ExpressionNode)
		if err != nil {
			return err
		}
		c.emit(code.OpPop)

	case *ast.InfixExpressionNode:
		if node.Operator == "<" {
			err := c.Compile(node.RightNode)
			if err != nil {
				return err
			}

			err = c.Compile(node.LeftNode)
			if err != nil {
				return err
			}
			c.emit(code.OpGreaterThan)
			return nil
		}

		err := c.Compile(node.LeftNode)
		if err != nil {
			return err
		}

		err = c.Compile(node.RightNode)
		if err != nil {
			return err
		}

		switch node.Operator {
		case "+":
			c.emit(code.OpAdd)
		case "-":
			c.emit(code.OpSub)
		case "*":
			c.emit(code.OpMul)
		case "/":
			c.emit(code.OpDiv)
		case ">":
			c.emit(code.OpGreaterThan)
		case "==":
			c.emit(code.OpEqual)
		case "!=":
			c.emit(code.OpNotEqual)
		default:
			return fmt.Errorf("unknown operator %s", node.Operator)
		}

	case *ast.IntegerLiteralNode:
		integer := &object.IntObject{Value: node.Value}
		c.emit(code.OpConstant, c.addConstant(integer))

	case *ast.BooleanNode:
		if node.Value {
			c.emit(code.OpTrue)
		} else {
			c.emit(code.OpFalse)
		}

	case *ast.PrefixExpressionNode:
		err := c.Compile(node.RightNode)
		if err != nil {
			return err
		}

		switch node.Operator {
		case "!":
			c.emit(code.OpBang)
		case "-":
			c.emit(code.OpMinus)
		default:
			return fmt.Errorf("unknown operator %s", node.Operator)
		}

	case *ast.IfExpressionNode:
		err := c.Compile(node.ConditionNode)
		if err != nil {
			return err
		}

		// Emit an `OpJumpNotTruthy` with a bogus value
		jumpNotTruthyPos := c.emit(code.OpJumpNotTruthy, 9999)

		err = c.Compile(node.ConsequenceNode)
		if err != nil {
			return err
		}

		if c.lastInstructionIs(code.OpPop) {
			c.removeLastPop()
		}

		// Emit an `OpJump` with a bogus value
		jumpPos := c.emit(code.OpJump, 9999)

		afterConsequencePos := len(c.currentInstructions())
		c.changeOperand(jumpNotTruthyPos, afterConsequencePos)

		if node.AlternativeNode == nil {
			c.emit(code.OpNull)
		} else {
			err := c.Compile(node.AlternativeNode)
			if err != nil {
				return err
			}

			if c.lastInstructionIs(code.OpPop) {
				c.removeLastPop()
			}
		}

		afterAlternativePos := len(c.currentInstructions())
		c.changeOperand(jumpPos, afterAlternativePos)

	case *ast.BlockStatementNode:
		for _, s := range node.StatementNodes {
			err := c.Compile(s)
			if err != nil {
				return err
			}
		}

	case *ast.LetStatementNode:
		err := c.Compile(node.ValueNode)
		if err != nil {
			return err
		}
		symbol := c.symbolTable.Define(node.NameNode.Value)
		if symbol.Scope == GlobalScope {
			c.emit(code.OpSetGlobal, symbol.Index)
		} else {
			c.emit(code.OpSetLocal, symbol.Index)
		}

	case *ast.IdentifierNode:
		symbol, ok := c.symbolTable.Resolve(node.Value)
		if !ok {
			return fmt.Errorf("undefined variable %s", node.Value)
		}

		c.loadSymbol(symbol)

	case *ast.StringLiteralNode:
		str := &object.StringObject{Value: node.Value}
		c.emit(code.OpConstant, c.addConstant(str))

	case *ast.ArrayLiteralNode:
		for _, el := range node.Elements {
			err := c.Compile(el)
			if err != nil {
				return err
			}
		}

		c.emit(code.OpArray, len(node.Elements))

	case *ast.HashLiteralNode:
		keys := []ast.ExpressionNode{}
		for k := range node.Pairs {
			keys = append(keys, k)
		}
		sort.Slice(keys, func(i, j int) bool {
			return keys[i].String() < keys[j].String()
		})

		for _, k := range keys {
			err := c.Compile(k)
			if err != nil {
				return err
			}
			err = c.Compile(node.Pairs[k])
			if err != nil {
				return err
			}
		}

		c.emit(code.OpHash, len(node.Pairs)*2)

	case *ast.IndexExpressionNode:
		err := c.Compile(node.Left)
		if err != nil {
			return err
		}

		err = c.Compile(node.Index)
		if err != nil {
			return err
		}

		c.emit(code.OpIndex)

	case *ast.FunctionLiteralNode:
		c.enterScope()

		for _, p := range node.ParamNodes {
			c.symbolTable.Define(p.Value)
		}

		err := c.Compile(node.BodyNode)
		if err != nil {
			return err
		}

		if c.lastInstructionIs(code.OpPop) {
			c.repaceLastPopWithReturn()
		}
		if !c.lastInstructionIs(code.OpReturnValue) {
			c.emit(code.OpReturn)
		}

		numLocals := c.symbolTable.counter
		instructions := c.leaveScope()

		compiledFunc := &object.CompiledFnObject{
			Instructions:  instructions,
			NumLocals:     numLocals,
			NumParameters: len(node.ParamNodes),
		}
		c.emit(code.OpConstant, c.addConstant(compiledFunc))

	case *ast.ReturnStatementNode:
		err := c.Compile(node.ReturnValueNode)
		if err != nil {
			return err
		}

		c.emit(code.OpReturnValue)

	case *ast.CallExpressionNode:
		err := c.Compile(node.FnNode)
		if err != nil {
			return err
		}

		for _, a := range node.ArgNodes {
			err := c.Compile(a)
			if err != nil {
				return err
			}
		}

		c.emit(code.OpCall, len(node.ArgNodes))

	}

	return nil
}
