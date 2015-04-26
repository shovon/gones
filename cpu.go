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
  // TODO: the instruction set is not just a set of bytes, but a bit more
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

// Gets the current content of the A register.
func (c* CPU) A() byte { return c.a }

// Gets the current content of the P register.
func (c* CPU) P() byte { return c.p }

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
  value := c.memory.GetUint8At(c.pc)
  c.pc++
  return value
}

// Gets the zero page address.
func (c* CPU) getZeroPageAddress() uint16 {
  c.cycles++
  return uint16(c.getFromImmediate())
}

// Gets the 8-bit value located at the zero page address.
func (c* CPU) getFromZeroPage() byte {
  return c.memory.GetUint8At(c.getZeroPageAddress());
}

// Gets the zero page,X address.
func (c* CPU) getZeroPageXAddress() uint16 {
  c.cycles += 2;
  return uint16(c.getFromImmediate() + c.x)
}

// Gets the 8-bit value located at the zero-page + x address.
func (c* CPU) getFromZeroPageX() byte {
  return c.memory.GetUint8At(c.getZeroPageXAddress());
}

// Gets the absolute address.
func (c* CPU) getAbsoluteAddress() uint16 {
  lsb := c.getFromImmediate()
  msb := c.getFromImmediate()
  address := (uint16(msb) << 8) & uint16(lsb)
  return address
}

// Gets the 8-bit value located at the absolute address.
func (c* CPU) getFromAbsolute() byte {
  return c.memory.GetUint8At(c.getAbsoluteAddress())
}

// This gets the absolute address with an offset.
func (c* CPU) getAbsoluteAddressWithOffset(offset byte, precompute bool) uint16 {
  lsb := c.getFromImmediate()
  msb := c.getFromImmediate()
  if (255 - offset < lsb || !precompute) {
    // This implies that page boundary has crossed.
    c.cycles++
  }
  address := ((uint16(msb) << 8) & uint16(lsb)) + uint16(offset)
  return address;
}

// Gets the Absolute,X address
func (c* CPU) getAbsoluteXAddress(precompute bool) uint16 {
  return c.getAbsoluteAddressWithOffset(c.x, precompute)
}

// Gets the 8-bit value located at the absolute + X address.
func (c* CPU) getFromAbsoluteX() byte {
  return c.memory.GetUint8At(c.getAbsoluteXAddress(true))
}

func (c* CPU) getAbsoluteYAddress(precompute bool) uint16 {
  return c.getAbsoluteAddressWithOffset(c.y, precompute)
}

// Gets the 8-bit value located at the absolute + Y address.
func (c* CPU) getFromAbsoluteY() byte {
  return c.memory.GetUint8At(c.getAbsoluteAddressWithOffset(c.y, true))
}

// Gets the Indirect,X address.
func (c* CPU) getIndirectIndexedAddress(precompute bool) uint16 {
  zeroPageAddress := c.getFromImmediate() + c.x
  lsb := c.memory.GetUint8At(uint16(zeroPageAddress))
  msb := c.memory.GetUint8At(uint16(zeroPageAddress + 1))
  if (!precompute || 255 - c.x < lsb) {
    c.cycles++
  }
  address := uint16(msb << 8) & uint16(lsb)
  return address;
}

// Gets the 8-bit value located at the indirect indexed address.
func (c* CPU) getFromIndirectIndexed() byte {
  return c.memory.GetUint8At(c.getIndirectIndexedAddress(true))
}

// Gets the (Indirect),Y address.
func (c* CPU) getIndexedIndirectAddress() uint16 {
  zeroPageAddress := c.getFromImmediate()
  lsb := c.memory.GetUint8At(uint16(zeroPageAddress))
  msb := c.memory.GetUint8At(uint16(zeroPageAddress + 1))
  address := uint16(msb << 8) & uint16(lsb) + uint16(c.y)
  return address
}

// Gets the 8-bit value located at the indexed indirect address.
func (c* CPU) getFromIndexedIndirect() byte {
  return c.memory.GetUint8At(c.getIndexedIndirectAddress())
}

// ADd with Carry
func (c* CPU) adc(value byte) {
  var carry byte = 0; if (c.C()) { carry = 1 }
  c.a = value + c.a + carry
}

// LoaD Accumulator
func (c* CPU) lda(value byte) {
  if value == 0 {
    c.SetZ(true)
  } else {
    c.SetZ(false)
  }
  if value & 0x80 != 0 {
    c.SetN(true)
  } else {
    c.SetN(false)
  }
  c.a = value
}

func (c *CPU) nop() {
  c.cycles++
}

func (c* CPU) sta(address uint16) {
  c.memory.SetUint8At(address, c.a)
}

// Simply runs the next instruction. Will write to registers and memory.
func (c* CPU) RunNextInstruction() error {
  switch c.getFromImmediate() {
  default: return errors.New("Opcode not supported")
  // ADC (ADd with Carry)
  case 0x69: c.adc(c.getFromImmediate())
  case 0x65: c.adc(c.getFromZeroPage())
  case 0x75: c.adc(c.getFromZeroPageX())
  case 0x6D: c.adc(c.getFromAbsolute())
  case 0x7D: c.adc(c.getFromAbsoluteX())
  case 0x79: c.adc(c.getFromAbsoluteY())
  case 0x61: c.adc(c.getFromIndirectIndexed())
  case 0x71: c.adc(c.getFromIndexedIndirect())

  // LDA (LoaD Accumulator)
  case 0xA9: c.lda(c.getFromImmediate())
  case 0xA5: c.lda(c.getFromZeroPage())
  case 0xB5: c.lda(c.getFromZeroPageX())
  case 0xAD: c.lda(c.getFromAbsolute())
  case 0xBD: c.lda(c.getFromAbsoluteX())
  case 0xB9: c.lda(c.getFromAbsoluteY())
  case 0xA1: c.lda(c.getFromIndirectIndexed())
  case 0xB1: c.lda(c.getFromIndexedIndirect())

  case 0xEA: c.nop()

  // STA (STore Accumulator)
  case 0x85: c.sta(c.getZeroPageAddress())
  case 0x95: c.sta(c.getZeroPageXAddress())
  case 0x8D: c.sta(c.getAbsoluteAddress())
  case 0x9D: c.sta(c.getAbsoluteXAddress(false))
  case 0x99: c.sta(c.getAbsoluteYAddress(false))
  case 0x81: c.sta(c.getIndexedIndirectAddress())
  case 0x90: c.sta(c.getIndirectIndexedAddress(false))

  // case 0x86: c.sta()
  }

  return nil
}

func (c* CPU) MovePCToResetVector() {
  c.pc = c.memory.GetUint16LEAt(0xFFFC)
}

// Starts the program in memory.
func (c* CPU) Run() int {
  c.MovePCToResetVector()

  for {
    c.RunNextInstruction()
    for c.cycles > 0 {
      c.cycles--
    }
  }

  return 0
}

// Converts a simple program to one that the 6502 can understand (that is, it
// resizes the size of the program to fit between memory locations 0x8000 and
// 0xFFFF, and adds the memory location where the program starts to the vector
// 0xFFFC and 0xFFFD)
func ConvertSimpleInstructions(instructions []byte) []byte {
  newInstructions := make([]byte, 0x8000, 0x8000)
  copy(newInstructions[0:], instructions)
  newInstructions[0x7FFC] = 0x00
  newInstructions[0x7FFD] = 0x80
  return newInstructions
}
