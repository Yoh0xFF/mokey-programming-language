package code

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type Opcode byte

type Definition struct {
	Name          string
	OperandWidths []int
}

type Instructions []byte

func (inst Instructions) String() string {
	var out bytes.Buffer

	i := 0
	for i < len(inst) {
		def, err := Lookup(inst[i])
		if err != nil {
			fmt.Fprintf(&out, "ERROR: %s\n", err)
			continue
		}

		operands, read := ReadOperands(def, inst[i+1:])

		fmt.Fprintf(&out, "%04d %s\n", i, inst.fmtInstruction(def, operands))

		i += 1 + read
	}

	return out.String()
}

func (inst Instructions) fmtInstruction(def *Definition, operands []int) string {
	operandCount := len(def.OperandWidths)

	if len(operands) != operandCount {
		return fmt.Sprintf("ERROR: operand len %d doesn't match defined %d\n",
			len(operands), operandCount)
	}

	switch operandCount {
	case 0:
		return def.Name
	case 1:
		return fmt.Sprintf("%s %d", def.Name, operands[0])
	}

	return fmt.Sprintf("ERROR: upnhandled operand count for %s\n", def.Name)
}

func Lookup(op byte) (*Definition, error) {
	def, ok := definitions[Opcode(op)]
	if !ok {
		return nil, fmt.Errorf("OpCode %d undefined", op)
	}

	return def, nil
}

func Make(op Opcode, operands ...int) []byte {
	def, ok := definitions[op]
	if !ok {
		return []byte{}
	}

	instructionLen := 1
	for _, w := range def.OperandWidths {
		instructionLen += w
	}

	instruction := make([]byte, instructionLen)
	instruction[0] = byte(op)

	offset := 1
	for i, o := range operands {
		w := def.OperandWidths[i]

		switch w {
		case 2:
			binary.BigEndian.PutUint16(instruction[offset:], uint16(o))
		case 1:
			instruction[offset] = byte(o)
		}

		offset += w
	}

	return instruction
}

func ReadOperands(def *Definition, inst Instructions) ([]int, int) {
	operands := make([]int, len(def.OperandWidths))
	offset := 0

	for i, w := range def.OperandWidths {
		switch w {
		case 2:
			operands[i] = int(ReadUint16(inst[offset:]))
		case 1:
			operands[i] = int(ReadUint8(inst[offset:]))
		}

		offset += w
	}

	return operands, offset
}

func ReadUint16(inst Instructions) uint16 {
	return binary.BigEndian.Uint16(inst)
}

func ReadUint8(inst Instructions) uint8 {
	return uint8(inst[0])
}
