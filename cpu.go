package main

import "errors"

const (
  // The masks that represents the flags
  C byte = 1 << iota
  Z
  I
  D
  B
  _
  V
  N
)

// Represents the 2A03/2A07 CPU architecture.
type CPU struct {
  // Registers
  pc uint16
  sp, a, x, y, p byte
  cycles int
  memory Memory
}

// Initializes a new CPU.
func CPUNew() CPU {
  return CPU{
    // TODO: check to see whether or not the initial values are correct.
    pc: 0,
    sp: 0xFF,
    a: 0,
    x: 0,
    y: 0,
    p: 0,
    cycles: 0,
    memory: Memory{},
  }
}

func (c* CPU) SetInstructions(instructions []byte) {
  // TODO: the instruction set is not actually a set of bytes, but a bit more
  // complicated than that.

  c.memory.SetInstructions(instructions)
}

// Gets the current CPU status flags.
func (c* CPU) status(flag byte) bool {
  return (c.p & flag) != 0
}

// Sets the current CPU status flags.
func (c* CPU) setStatus(flag byte, status bool) {
  if status {
    c.p = c.p | flag
  } else {
    c.p = c.p & ^flag
  }
}

// Gets the current value of the C flag
func (c* CPU) C() bool { return c.status(C) }

// Sets the C flag
func (c* CPU) SetC(status bool) { c.setStatus(C, status) }

// Gets the value of the current Z flag
func (c* CPU) Z() bool { return c.status(Z) }

// Sets the Z flag
func (c* CPU) SetZ(status bool) { c.setStatus(Z, status) }

// Gets the current value of the I flag
func (c* CPU) I() bool { return c.status(I) }

// Sets the I flag
func (c* CPU) SetI(status bool) { c.status(I) }

// Gets the current value of the D flag
func (c* CPU) D() bool { return c.status(D) }

// Sets the D flag
func (c* CPU) SetD(status bool) { c.setStatus(D, status) }

// Gets the current value of the B flag
func (c* CPU) B() bool { return c.status(B) }

// Sets the B flag
func (c* CPU) SetB(status bool) { c.setStatus(B, status) }

// Gets the current value of the V flag
func (c* CPU) V() bool { return c.status(V) }

// Sets the V flag
func (c* CPU) SetV(status bool) { c.setStatus(V, status) }

// Gets the current value of the N flag
func (c* CPU) N() bool { return c.status(N) }

// Sets the N flag
func (c* CPU) SetN(status bool) { c.setStatus(N, status) }

// Gets the 8-bit value located where the program counter is pointing to.
func (c* CPU) getFromImmediate() byte {
  c.cycles++
  value := c.memory.GetInt8At(c.pc)
  c.pc++
  return value
}

// Gets the 8-bit value located at the zero page address.
func (c* CPU) getFromZeroPage() byte {
  c.cycles++
  return c.memory.GetInt8At(uint16(c.getFromImmediate()));
}

// Gets the 8-bit value located at the zero-page + x address.
func (c* CPU) getFromZeroPageX() byte {
  c.cycles += 2;
  return c.memory.GetInt8At(uint16(c.getFromImmediate() + c.x));
}

// Gets the 8-bit value located at the absolute address.
func (c* CPU) getFromAbsolute() byte {
  lsb := c.getFromImmediate()
  msb := c.getFromImmediate()
  address := (uint16(msb) << 8) & uint16(lsb)
  return c.memory.GetInt8At(address)
}

// This gets the absolute address with an offset.
func (c* CPU) getAbsoluteAddress(offset byte) uint16 {
  lsb := c.getFromImmediate()
  msb := c.getFromImmediate()
  if (255 - offset < lsb) {
    // This implies that page boundary has crossed.
    c.cycles++
  }
  address := ((uint16(msb) << 8) & uint16(lsb)) + uint16(offset)
  return address;
}

// Gets the 8-bit value located at the absolute + X address.
func (c* CPU) getFromAbsoluteX() byte {
  return c.memory.GetInt8At(c.getAbsoluteAddress(c.x))
}

// Gets the 8-bit value located at the absolute + Y address.
func (c* CPU) getFromAbsoluteY() byte {
  return c.memory.GetInt8At(c.getAbsoluteAddress(c.y))
}

// Gets the 8-bit value located at the indirect indexed address.
func (c* CPU) getFromIndirectIndexed() byte {
  zeroPageAddress := c.getFromImmediate() + c.x
  lsb := c.memory.GetInt8At(uint16(zeroPageAddress))
  msb := c.memory.GetInt8At(uint16(zeroPageAddress + 1))
  
  address := uint16(msb << 8) & uint16(lsb)

  return c.memory.GetInt8At(address)
}

// Gets the 8-bit value located at the indexed indirect address.
func (c* CPU) getFromIndexedIndirect() byte {
  zeroPageAddress := c.getFromImmediate()
  lsb := c.memory.GetInt8At(uint16(zeroPageAddress))
  msb := c.memory.GetInt8At(uint16(zeroPageAddress + 1))
  if (255 - c.y < lsb) {
    c.cycles++
  }
  address := uint16(msb << 8) & uint16(lsb) + uint16(c.y)

  return c.memory.GetInt8At(address)
}

func (c* CPU) adc(value byte) {
  var carry byte = 0; if (c.C()) { carry = 1 }
  c.a = value + c.a + carry
}

// Simply runs the next instruction. Will write to registers and memory.
func (c* CPU) RunNextInstruction() error {
  switch c.getFromImmediate() {
  default: return errors.New("Opcode not supported")
  // ADC (ADd with Carry)
  case 0x69:
    value := c.getFromImmediate()
    c.adc(value)
  case 0x65:
    value := c.getFromZeroPage()
    c.adc(value)
  case 0x75:
    value := c.getFromZeroPageX()
    c.adc(value)
  case 0x6D:
    value := c.getFromAbsolute()
    c.adc(value)
  case 0x7D:
    value := c.getFromAbsoluteX()
    c.adc(value)
  case 0x79:
    value := c.getFromAbsoluteY()
    c.adc(value)
  case 0x61:
    value := c.getFromIndirectIndexed()
    c.adc(value)
  case 0x71:
    value := c.getFromIndexedIndirect()
    c.adc(value)
  }

  return nil
}

func (c* CPU) Run() int {
  c.pc = c.memory.GetUint16LEAt(0xFFFC)

  for {
    c.RunNextInstruction()
    for c.cycles > 0 {
      c.cycles--
    }
  }

  return 0
}
