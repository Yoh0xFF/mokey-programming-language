package vm

import (
	"monkey/code"
	"monkey/object"
)

type Frame struct {
	fn          *object.CompiledFnObject
	ip          int
	basePointer int
}

func NewFrame(fn *object.CompiledFnObject, basePointer int) *Frame {
	frame := &Frame{
		fn:          fn,
		ip:          -1,
		basePointer: basePointer,
	}

	return frame
}

func (f *Frame) Instructions() code.Instructions {
	return f.fn.Instructions
}
