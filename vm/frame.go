package vm

import (
	"monkey/code"
	"monkey/object"
)

type Frame struct {
	cl          *object.ClosureObject
	ip          int
	basePointer int
}

func NewFrame(cl *object.ClosureObject, basePointer int) *Frame {
	frame := &Frame{
		cl:          cl,
		ip:          -1,
		basePointer: basePointer,
	}

	return frame
}

func (f *Frame) Instructions() code.Instructions {
	return f.cl.Fn.Instructions
}
