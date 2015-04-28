package main

import "testing"
import "log"

func testStatus(t *testing.T, flag byte) {
  cpu := CPUNew()
  if (cpu.status(flag)) {
    t.Fail()
  }

  cpu.setStatus(flag, true)

  if (!cpu.status(flag)) {
    t.Fail()
  }

  cpu.setStatus(flag, false)

  if (cpu.status(flag)) {
    t.Fail()
  }
}

func TestStatus(t *testing.T) {
  testStatus(t, C)
  testStatus(t, Z)
  testStatus(t, I)
  testStatus(t, D)
  testStatus(t, B)
  testStatus(t, V)
  testStatus(t, N)
}

func initCPUWithBasicInstructions(instructions []byte) *CPU {
  newInstructions := ConvertSimpleInstructions(instructions)
  cpu := CPUNew()
  cpu.SetInstructions(newInstructions)
  cpu.MovePCToResetVector()
  return cpu
}

func TestConvertSimpleInstructions(t *testing.T) {
  instructions := ConvertSimpleInstructions([]byte{
    0xEA,
  })
  if len(instructions) != 0x8000 {
    log.Printf("Was expecting instructions to be of size 0x8000 but got size %X", len(instructions))
    t.Fail()
  }
  resetVector := GetUint16LEAt(instructions, 0x7FFC)
  if resetVector != 0x8000 {
    log.Printf("Was expecting 0x8000 at location 0x7FFC, but got %X", resetVector)
    t.Fail()
  }
}

func TestSetInstructions(t *testing.T) {
  instructions := ConvertSimpleInstructions([]byte{
    0xEA,
  })
  cpu := CPUNew()
  cpu.SetInstructions(instructions)
  if cpu.memory.GetUint8At(0x8000) != 0xEA {
    log.Printf("Expecting program start to be 0xEA")
    t.Fail()
  }

  resetVector := cpu.memory.GetUint16LEAt(0xFFFC)
  if resetVector != 0x8000 {
    log.Printf("Expecting value at memory location 0xFFFC to be 0x8000, but got %d", resetVector)
    t.Fail()
  }
}

func TestMovePCToResetVector(t *testing.T) {
  instructions := ConvertSimpleInstructions([]byte{
    0xEA,
  })
  cpu := CPUNew()
  cpu.SetInstructions(instructions)
  cpu.MovePCToResetVector()
  if cpu.pc != 0x8000 {
    log.Printf("Expecting program counter to point to location 0x8000, but it's pointing at %d", cpu.pc)
    t.Fail()
  }
}

func TestAdc(t *testing.T) {
  
}

func TestLda(t *testing.T) {
  instructions := ConvertSimpleInstructions([]byte{
    // Immediate
    0xA9, 3,
    0xA9, 128,
    0xA9, 0,

    // Zero Page
    // Set-up
    0xA9, 42,
    0x85, 0x24,
    0xA9, 0,
    // Instruction under test.
    0xA5, 0x24,
  })
  cpu := CPUNew()
  cpu.SetInstructions(instructions)
  cpu.MovePCToResetVector()

  cpu.RunNextInstruction()
  if cpu.cycles != 2 {
    log.Printf("LDA immediate should have expanded two clock cycles")
    t.Fail()
  }
  if cpu.N() {
    log.Printf("The number in accumulator is not signed as negative, but the N flag has been set to true")
    t.Fail()
  }
  if cpu.Z() {
    log.Printf("The number in accumulator is not zero, but the Z flag has been set to true")
    t.Fail()
  }
  if cpu.A() != 3 {
    log.Printf("Expecting 3, but got %d", cpu.A)
    t.Fail()
  }
  cpu.cycles = 0

  cpu.RunNextInstruction()
  if cpu.cycles != 2 { t.Fail() }
  if !cpu.N() {
    log.Printf("The number in accumulator is signed as negative, but the N flag is set to false")
    t.Fail()
  }
  if cpu.Z() { t.Fail() }
  if cpu.A() != 128 { t.Fail() }
  cpu.cycles = 0

  cpu.RunNextInstruction()
  if cpu.cycles != 2 { t.Fail() }
  if cpu.N() { t.Fail() }
  if !cpu.Z() { t.Fail() }
  if cpu.A() != 0 { t.Fail() }
  cpu.cycles = 0

  cpu.RunNextInstruction(); cpu.cycles = 0
  cpu.RunNextInstruction(); cpu.cycles = 0
  cpu.RunNextInstruction();
  if cpu.A() != 0 { t.Fail() } // Reset to test another operation, below
  cpu.cycles = 0
  cpu.RunNextInstruction()
  if cpu.A() != 42 { t.Fail() }
  cpu.cycles = 0
}

func TestLdx(t *testing.T) {
  instructions := ConvertSimpleInstructions([]byte{
    // Immediate
    0xA2, 3,
    0xA2, 128,
    0xA2, 0,

    // Zero Page
    // Set-up
    0xA9, 42,
  })
  cpu := CPUNew()
  cpu.SetInstructions(instructions)
  cpu.MovePCToResetVector()

  var oldA byte

  oldA = cpu.A()

  cpu.RunNextInstruction()
  if cpu.cycles != 2 { t.Fail() }
  if cpu.A() != oldA { t.Fail() }
  if cpu.X() != 3 { t.Fail() }
  if cpu.Z() { t.Fail() }
  if cpu.N() { t.Fail() }
  cpu.cycles = 0

  cpu.RunNextInstruction()
  if cpu.cycles != 2 { t.Fail() }
  if cpu.A() != oldA { t.Fail() }
  if cpu.X() != 128 { t.Fail() }
  if cpu.Z() { t.Fail() }
  if !cpu.N() { t.Fail() }
  cpu.cycles = 0

  cpu.RunNextInstruction()
  if cpu.cycles != 2 { t.Fail() }
  if cpu.A() != oldA { t.Fail() }
  if cpu.X() != 0 { t.Fail() }
  if !cpu.Z() { t.Fail() }
  if cpu.N() { t.Fail() }
  cpu.cycles = 0
}

func TestNop(t *testing.T) {
  instructions := ConvertSimpleInstructions([]byte{
    0xEA,
  })
  cpu := CPUNew()
  cpu.SetInstructions(instructions)
  cpu.MovePCToResetVector()

  cpu.RunNextInstruction();
  if cpu.cycles != 2 {
    log.Printf("Expecting 2 CPU cycles, but got %d", cpu.cycles)
    t.Fail()
  }
  cpu.cycles = 0
}

func TestSta(t *testing.T) {
  instructions := ConvertSimpleInstructions([]byte{
    0xA9, 42,
    0x85, 0x24,
  })
  cpu := CPUNew()
  cpu.SetInstructions(instructions)
  cpu.MovePCToResetVector()
  
  cpu.RunNextInstruction(); cpu.cycles = 0
  previousP := cpu.P()
  cpu.RunNextInstruction()
  if cpu.cycles != 3 { t.Fail() }
  if cpu.P() != previousP { t.Fail() }
  if cpu.memory.GetUint8At(0x24) != 42 { t.Fail() }
  if cpu.A() != 42 { t.Fail() }
  cpu.cycles = 0
}

func TestStx(t *testing.T) {
  instructions := ConvertSimpleInstructions([]byte{
    // Zero Page
    // Set-up
    0xA2, 42,   // LDX #42 ; Load 42 onto register x
    0x86, 0x24, // STX $0  ; Store content of register X at zero page address $0
  })
  cpu := CPUNew()
  cpu.SetInstructions(instructions)
  cpu.MovePCToResetVector()

  cpu.RunNextInstruction(); cpu.cycles = 0
  if cpu.X() != 42 {
    log.Printf("Expecting X register to be 42, but was actually %d", cpu.X())
    t.Fail()
  }
  cpu.RunNextInstruction()
  if cpu.cycles != 3 {
    log.Printf("Expecting CPU cycles to be 3, but got %d", cpu.cycles)
    t.Fail()
  }
  if cpu.memory.GetUint8At(0x24) != 42 {
    log.Printf(
      "Expecting value at location 0x24 to be 42, but got %d",
      cpu.memory.GetUint8At(0x24),
    )
    t.Fail()
  }
  if cpu.X() != 42 { t.Fail() }
  cpu.cycles = 0
}