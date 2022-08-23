package vm

import (
	"fmt"
	"monkey/code"
	"monkey/compiler"
	"monkey/object"
)

const StackSize = 2048
const GlobalsSize = 65536
const MaxFrames = 1024

var True = &object.BoolObject{Value: true}
var False = &object.BoolObject{Value: false}
var Null = &object.NullObject{}

type VM struct {
	constants []object.Object
	globals   []object.Object

	stack []object.Object
	sp    int // Always points to the next value. Top of the stack is stack[sp - 1]

	frames      []*Frame
	framesIndex int
}

func New(bytecode *compiler.Bytecode) *VM {
	mainFn := &object.CompiledFnObject{Instructions: bytecode.Instructions}
	mainClosure := &object.ClosureObject{Fn: mainFn}
	mainFrame := NewFrame(mainClosure, 0)

	frames := make([]*Frame, MaxFrames)
	frames[0] = mainFrame

	return &VM{
		constants: bytecode.Constants,
		globals:   make([]object.Object, GlobalsSize),

		stack: make([]object.Object, StackSize),
		sp:    0,

		frames:      frames,
		framesIndex: 1,
	}
}

func NewWithGlobalsStore(bytecode *compiler.Bytecode, globals []object.Object) *VM {
	vm := New(bytecode)
	vm.globals = globals
	return vm
}

func (vm *VM) LastPoppedStackElem() object.Object {
	return vm.stack[vm.sp]
}

func (vm *VM) push(o object.Object) error {
	if vm.sp >= StackSize {
		return fmt.Errorf("stack overflow")
	}

	vm.stack[vm.sp] = o
	vm.sp++

	return nil
}

func (vm *VM) pop() object.Object {
	o := vm.stack[vm.sp-1]
	vm.sp--
	return o
}

func (vm *VM) executeBinaryOperation(op code.Opcode) error {
	right := vm.pop()
	left := vm.pop()

	leftType := left.Type()
	rightType := right.Type()

	switch {
	case leftType == object.INT_OBJ && rightType == object.INT_OBJ:
		return vm.executeBinaryIntegerOperation(op, left, right)
	case leftType == object.STRING_OBJ && rightType == object.STRING_OBJ:
		return vm.executeBinaryStringOperation(op, left, right)
	default:
		return fmt.Errorf("unsupported types for binary operation: %s %s",
			leftType, rightType)
	}
}

func (vm *VM) executeBinaryIntegerOperation(
	op code.Opcode,
	left, right object.Object,
) error {
	leftValue := left.(*object.IntObject).Value
	rightValue := right.(*object.IntObject).Value

	var result int64

	switch op {
	case code.OpAdd:
		result = leftValue + rightValue
	case code.OpSub:
		result = leftValue - rightValue
	case code.OpMul:
		result = leftValue * rightValue
	case code.OpDiv:
		result = leftValue / rightValue
	default:
		return fmt.Errorf("unknown integer operator: %d", op)
	}

	return vm.push(&object.IntObject{Value: result})
}

func (vm *VM) executeComparison(op code.Opcode) error {
	right := vm.pop()
	left := vm.pop()

	if left.Type() == object.INT_OBJ && right.Type() == object.INT_OBJ {
		return vm.executeIntegerComparison(op, left, right)
	}

	switch op {
	case code.OpEqual:
		return vm.push(nativeBoolToBooleanObject(right == left))
	case code.OpNotEqual:
		return vm.push(nativeBoolToBooleanObject(right != left))
	default:
		return fmt.Errorf("unknown operator: %d (%s %s)",
			op, left.Type(), right.Type())
	}
}

func (vm *VM) executeIntegerComparison(
	op code.Opcode,
	left, right object.Object,
) error {
	leftValue := left.(*object.IntObject).Value
	rightValue := right.(*object.IntObject).Value

	switch op {
	case code.OpEqual:
		return vm.push(nativeBoolToBooleanObject(rightValue == leftValue))
	case code.OpNotEqual:
		return vm.push(nativeBoolToBooleanObject(rightValue != leftValue))
	case code.OpGreaterThan:
		return vm.push(nativeBoolToBooleanObject(leftValue > rightValue))
	default:
		return fmt.Errorf("unknown operator: %d", op)
	}
}

func (vm *VM) executeBangOperator() error {
	operand := vm.pop()

	switch operand {
	case True:
		return vm.push(False)
	case False:
		return vm.push(True)
	case Null:
		return vm.push(True)
	default:
		return vm.push(False)
	}
}

func (vm *VM) executeMinusOperator() error {
	operand := vm.pop()

	if operand.Type() != object.INT_OBJ {
		return fmt.Errorf("unsupported type for negation: %s", operand.Type())
	}

	value := operand.(*object.IntObject).Value
	return vm.push(&object.IntObject{Value: -value})
}

func (vm *VM) executeBinaryStringOperation(
	op code.Opcode,
	left, right object.Object,
) error {
	if op != code.OpAdd {
		return fmt.Errorf("unknown string operator: %d", op)
	}

	leftValue := left.(*object.StringObject).Value
	rightValue := right.(*object.StringObject).Value

	return vm.push(&object.StringObject{Value: leftValue + rightValue})
}

func (vm *VM) buildArray(startIndex, endIndex int) object.Object {
	elements := make([]object.Object, endIndex-startIndex)

	for i := startIndex; i < endIndex; i++ {
		elements[i-startIndex] = vm.stack[i]
	}

	return &object.ArrayObject{Elements: elements}
}

func (vm *VM) buildHash(startIndex, endIndex int) (object.Object, error) {
	hashPairs := make(map[object.HashKey]object.HashPair)

	for i := startIndex; i < endIndex; i += 2 {
		key := vm.stack[i]
		value := vm.stack[i+1]

		pair := object.HashPair{Key: key, Value: value}

		hashKey, ok := key.(object.Hashable)
		if !ok {
			return nil, fmt.Errorf("unusable as hash key: %s", key.Type())
		}

		hashPairs[hashKey.HashKey()] = pair
	}

	return &object.HashObject{Pairs: hashPairs}, nil
}

func (vm *VM) executeIndexExpression(left, index object.Object) error {
	switch {
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.INT_OBJ:
		return vm.executeArrayIndex(left, index)
	case left.Type() == object.HASH_OBJ:
		return vm.executeHashIndex(left, index)
	default:
		return fmt.Errorf("index operator not supported: %s", left.Type())
	}
}

func (vm *VM) executeArrayIndex(array, index object.Object) error {
	arrayObject := array.(*object.ArrayObject)
	i := index.(*object.IntObject).Value
	max := int64(len(arrayObject.Elements) - 1)

	if i < 0 || i > max {
		return vm.push(Null)
	}

	return vm.push(arrayObject.Elements[i])
}

func (vm *VM) executeHashIndex(hash, index object.Object) error {
	hashObject := hash.(*object.HashObject)

	key, ok := index.(object.Hashable)
	if !ok {
		return fmt.Errorf("unusable as hash key: %s", index.Type())
	}

	pair, ok := hashObject.Pairs[key.HashKey()]
	if !ok {
		return vm.push(Null)
	}

	return vm.push(pair.Value)
}

func (vm *VM) currentFrame() *Frame {
	return vm.frames[vm.framesIndex-1]
}

func (vm *VM) pushFrame(f *Frame) {
	vm.frames[vm.framesIndex] = f
	vm.framesIndex++
}

func (vm *VM) popFrame() *Frame {
	vm.framesIndex--
	return vm.frames[vm.framesIndex]
}

func (vm *VM) executeCall(numArgs int) error {
	callee := vm.stack[vm.sp-1-numArgs]
	switch callee := callee.(type) {
	case *object.ClosureObject:
		return vm.callClosure(callee, numArgs)
	case *object.BuiltinObject:
		return vm.callBuiltin(callee, numArgs)
	default:
		return fmt.Errorf("calling non-function and non-built-in")
	}
}

func (vm *VM) callClosure(cl *object.ClosureObject, numArgs int) error {
	if numArgs != cl.Fn.NumParameters {
		return fmt.Errorf("wrong number of arguments: want=%d, got=%d",
			cl.Fn.NumParameters, numArgs)
	}

	frame := NewFrame(cl, vm.sp-numArgs)
	vm.pushFrame(frame)

	vm.sp = frame.basePointer + cl.Fn.NumLocals

	return nil
}

func (vm *VM) callBuiltin(builtin *object.BuiltinObject, numArgs int) error {
	args := vm.stack[vm.sp-numArgs : vm.sp]

	result := builtin.Fn(args...)
	vm.sp = vm.sp - numArgs - 1

	if result != nil {
		vm.push(result)
	} else {
		vm.push(Null)
	}

	return nil
}

func (vm *VM) pushClosure(constIndex, numFree int) error {
	constant := vm.constants[constIndex]
	fn, ok := constant.(*object.CompiledFnObject)
	if !ok {
		return fmt.Errorf("not a function: %+v", constant)
	}

	free := make([]object.Object, numFree)
	for i := 0; i < numFree; i++ {
		free[i] = vm.stack[vm.sp-numFree+i]
	}
	vm.sp = vm.sp - numFree

	closure := &object.ClosureObject{Fn: fn, Free: free}
	return vm.push(closure)
}

func nativeBoolToBooleanObject(input bool) *object.BoolObject {
	if input {
		return True
	}
	return False
}

func isTruthy(obj object.Object) bool {
	switch obj := obj.(type) {

	case *object.BoolObject:
		return obj.Value

	case *object.NullObject:
		return false

	default:
		return true
	}
}
