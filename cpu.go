package main

import "errors"

const (
  // The masks that represent the flags
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
func CPUNew() *CPU {
  return &CPU{
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

// Gets the curretn content of the X register.
func (c* CPU) X() byte { return c.x }

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
//
// Adds a CPU cycle, and advances the program counter by one.
func (c* CPU) getFromImmediate() byte {
  c.cycles++
  value := c.memory.GetUint8At(c.pc)
  c.pc++
  return value
}

// Gets the zero page address.
// 
// Adds two CPU cycles, and advances the program counter by one.
func (c* CPU) getZeroPageAddress() uint16 {
  c.cycles++
  return uint16(c.getFromImmediate())
}

// Gets the 8-bit value located at the zero page address.
//
// Adds two CPU cycles, and advances the program counter by one.
func (c* CPU) getFromZeroPage() byte {
  address := c.getZeroPageAddress()
  value := c.memory.GetUint8At(address)
  return value;
}

// Gets the Zero Page,X address.
//
// Adds three CPU cycles, and advances the program counter by one.
func (c* CPU) getZeroPageXAddress() uint16 {
  c.cycles += 2;
  return uint16(c.getFromImmediate() + c.x)
}

// Gets the 8-bit value located at the zero-page + x address.
func (c* CPU) getFromZeroPageX() byte {
  return c.memory.GetUint8At(c.getZeroPageXAddress())
}

// Gets the Zero Page,Y address.
//
// Adds four CPU cycles, and advances the program counter by one.
func (c* CPU) getZeroPageYAddress() uint16 {
  c.cycles += 3;
  return uint16(c.getFromImmediate() + c.y)
}

// Gets the 8-but value located at the zero-age + y address.
func (c* CPU) getFromZereoPageY() byte {
  return c.memory.GetUint8At(c.getZeroPageYAddress())
}

// Gets the absolute address.
//
// Adds 
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

func isNegative(value byte) bool {
  return value & 0x80 != 0
}

// ADd with Carry
func (c* CPU) adc(value byte) {
  var carry byte = 0; if (c.C()) { carry = 1 }
  c.a = value + c.a + carry
}

// LoaD Accumulator
func (c* CPU) lda(value byte) {
  c.SetZ(value == 0)
  c.SetN(isNegative(value))
  c.a = value
}

// LoaD X register
func (c* CPU) ldx(value byte) {
  c.SetZ(value == 0)
  c.SetN(isNegative(value))
  c.x = value
}

// N OPeration
func (c *CPU) nop() {
  c.cycles++
}

// STore Accumulator
func (c* CPU) sta(address uint16) {
  c.memory.SetUint8At(address, c.A())
}

// STore X register
func (c* CPU) stx(address uint16) {
  c.memory.SetUint8At(address, c.x)
}

// Simply runs the next instruction. Will write to registers and memory.
func (c* CPU) RunNextInstruction() error {
  switch c.getFromImmediate() {

  default: return errors.New("Opcode not supported")

  // ADC (ADd with Carry)
  // TODO: test this
  case 0x69: c.adc(c.getFromImmediate())
  // TODO: test this
  case 0x65: c.adc(c.getFromZeroPage())
  // TODO: test this
  case 0x75: c.adc(c.getFromZeroPageX())
  // TODO: test this
  case 0x6D: c.adc(c.getFromAbsolute())
  // TODO: test this
  case 0x7D: c.adc(c.getFromAbsoluteX())
  // TODO: test this
  case 0x79: c.adc(c.getFromAbsoluteY())
  // TODO: test this
  case 0x61: c.adc(c.getFromIndirectIndexed())
  // TODO: test this
  case 0x71: c.adc(c.getFromIndexedIndirect())

  // AND (Logical AND)
  // TODO: implement AND

  // ASL (Arithmetic Shift Left)
  // TODO: implement ASL

  // BCC (Branch if Carry Clear)
  // TODO: implement BCC

  // BCS (Branch if Carry Set)
  // TODO: implement BCS

  // BEQ (Branch if EQual)
  // TODO: implement BEQ

  // BIT (BIT test)
  // TODO: implement BIT

  // BMI (Branch if MInus)
  // TODO: implement BMI

  // BNE (Branch if Not Equal)
  // TODO: impelement BNE

  // BPL (Branch if positive (PLus))
  // TODO: implement BPL

  // BRK (force interrupt (BReaK))
  // TODO: implement BRK

  // BVC (Branch if oVerflow Clear)
  // TODO: implement BVC

  // BVS (Branch if oVerflow Set)
  // TODO: implement BVS

  // CLC (CLear Carry flag)
  // TODO: implement CLC

  // CLD (CLear Decimal mode)
  // TODO: implement CLD

  // LDA (LoaD Accumulator)
  case 0xA9: c.lda(c.getFromImmediate())
  case 0xA5: c.lda(c.getFromZeroPage())
  // TODO: test this
  case 0xB5: c.lda(c.getFromZeroPageX())
  // TODO: test this
  case 0xAD: c.lda(c.getFromAbsolute())
  // TODO: test this
  case 0xBD: c.lda(c.getFromAbsoluteX())
  // TODO: test this
  case 0xB9: c.lda(c.getFromAbsoluteY())
  // TODO: test this
  case 0xA1: c.lda(c.getFromIndirectIndexed())
  // TODO: test this
  case 0xB1: c.lda(c.getFromIndexedIndirect())

  // LDX (LoaD X Register)
  case 0xA2: c.ldx(c.getFromImmediate())
  case 0xA6: c.ldx(c.getFromZeroPage())
  // TODO: test this
  case 0xB6: c.ldx(c.getFromZeroPageY())
  // TODO: test this
  case 0xAE: c.ldx(c.getFromAbsolute())
  // TODO: test this
  case 0xBE: c.ldx(c.getFromAbsoluteY())

  // LDY (LoaD Y Register)
  // TODO: implement LDY

  // NOP (NO oPeration)
  case 0xEA: c.nop()

  // STA (STore Accumulator)
  // TODO: test this
  case 0x85: c.sta(c.getZeroPageAddress())
  // TODO: test this
  case 0x95: c.sta(c.getZeroPageXAddress())
  // TODO: test this
  case 0x8D: c.sta(c.getAbsoluteAddress())
  // TODO: test this
  case 0x9D: c.sta(c.getAbsoluteXAddress(false))
  // TODO: test this
  case 0x99: c.sta(c.getAbsoluteYAddress(false))
  // TODO: test this
  case 0x81: c.sta(c.getIndexedIndirectAddress())
  // TODO: test this
  case 0x90: c.sta(c.getIndirectIndexedAddress(false))

  // STX (STore X register)
  case 0x86: c.stx(c.getZeroPageAddress())
  }

  return nil
}

// Have the program counter point to the location represented by the 16-bit LE
// values located at addresses 0xFFFC
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
