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

////////////////////////////////////////////////////////////////////////////////

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
    0xA9, 42,   // LDA #42
    0x85, 0x24, // STA $24
    // Our operation under test
    0xA6, 0x24, // LDX $24
  })
  cpu := CPUNew()
  cpu.SetInstructions(instructions)
  cpu.MovePCToResetVector()

  var oldA byte

  oldA = cpu.A()

  // LDX #3
  cpu.RunNextInstruction()
  if cpu.cycles != 2 { t.Fail() }
  if cpu.A() != oldA { t.Fail() }
  if cpu.X() != 3 { t.Fail() }
  if cpu.Z() { t.Fail() }
  if cpu.N() { t.Fail() }
  cpu.cycles = 0

  // LDX #128
  cpu.RunNextInstruction()
  if cpu.cycles != 2 { t.Fail() }
  if cpu.A() != oldA { t.Fail() }
  if cpu.X() != 128 { t.Fail() }
  if cpu.Z() { t.Fail() }
  if !cpu.N() { t.Fail() }
  cpu.cycles = 0

  // LDX #0
  cpu.RunNextInstruction()
  if cpu.cycles != 2 { t.Fail() }
  if cpu.A() != oldA { t.Fail() }
  if cpu.X() != 0 { t.Fail() }
  if !cpu.Z() { t.Fail() }
  if cpu.N() { t.Fail() }
  cpu.cycles = 0

  if cpu.memory.GetUint8At(cpu.pc) != 0xA9 {
    log.Printf("Expecting immediate LDA, but got a different instruction")
    t.Fail()
  }

  // LDA #42
  if cpu.RunNextInstruction() != nil {
    log.Printf("An error occurred")
    t.Fail()
  }
  cpu.cycles = 0; oldA = cpu.A()
  if cpu.A() != 42 {
    log.Printf("Expecting accumulator to have value 42 but got %d", cpu.A())
    t.Fail()
  }

  // STA $24
  cpu.RunNextInstruction(); cpu.cycles = 0

  // LDX $24
  cpu.RunNextInstruction();
  if cpu.cycles != 3 {
    log.Printf("Expecting 3 CPU cycles but got %d", cpu.cycles)
    t.Fail()
  }
  if cpu.A() != oldA {
    log.Printf("Expecting accumulator to be %d but got %d", oldA, cpu.A())
    t.Fail()
  }
  if cpu.X() != 42 {
    log.Printf("Expecting register X to be 42 but got %d", cpu.X())
    t.Fail()
  }
  cpu.cycles = 0
}

func TestLdy(t *testing.T) {
  instructions := ConvertSimpleInstructions([]byte{
    // Immediate
    0xA0, 3,
    0xA0, 128,
    0xA0, 0,

    // Zero Page
    // Set-up
    0xA9, 42,   // LDA #42
    0x85, 0x24, // STA $24
    // Our operation under test
    0xA4, 0x24, // LDY $24

    // Zero Page,X
    // Set-up
    0xA9, 24,   // LDA #24
    0xA2, 2,    // LDX #2
    0x95, 0x42, // STA $42,X
    // Our operation under test
    0xB4, 0x42,
  })
  cpu := CPUNew()
  cpu.SetInstructions(instructions)
  cpu.MovePCToResetVector()

  var oldA byte

  oldA = cpu.A()

  // LDY #3
  if cpu.RunNextInstruction() != nil {
    log.Printf("An error occurred")
    t.Fail()
  }
  if cpu.cycles != 2 { t.Fail() }
  if cpu.A() != oldA { t.Fail() }
  if cpu.Y() != 3 { t.Fail() }
  if cpu.Z() { t.Fail() }
  if cpu.N() { t.Fail() }
  cpu.cycles = 0

  // LDY #128
  if cpu.RunNextInstruction() != nil {
    log.Printf("An error occurred")
    t.Fail()
  }
  if cpu.cycles != 2 {
    log.Printf("Expecting cycles to be %d, but got %d", 2, cpu.cycles)
    t.Fail()
  }
  if cpu.A() != oldA {
    log.Printf("Expecting accumulator to be %d, but got %d", oldA, cpu.A())
    t.Fail()
  }
  if cpu.Y() != 128 {
    log.Printf("Expecting register Y to be %d, but got %d", 128, cpu.Y())
    t.Fail()
  }
  if cpu.Z() { t.Fail() }
  if !cpu.N() { t.Fail() }
  cpu.cycles = 0

  // LDY #0
  if cpu.RunNextInstruction() != nil {
    log.Printf("An error occurred")
  }
  if cpu.cycles != 2 {
    log.Printf("Expecting CPU cycles count to be %d, but got cpu.cycles", 2, cpu.cycles)
    t.Fail()
  }
  if cpu.A() != oldA {
    log.Printf("Expecting accumulator to have value %d, but got %d", oldA, cpu.A())
    t.Fail()
  }
  if cpu.Y() != 0 {
    log.Printf("Expecting register Y to have value %d but got %d", 0, cpu.Y())
    t.Fail()
  }
  if !cpu.Z() {
    log.Printf("Expecting flag Z to be set")
    t.Fail()
  }
  if cpu.N() {
    log.Printf("Expecting flag N to be clear")
    t.Fail()
  }
  cpu.cycles = 0

  if cpu.memory.GetUint8At(cpu.pc) != 0xA9 {
    log.Printf("Expecting immediate LDA, but got a different instruction")
    t.Fail()
  }

  // LDA #42
  cpu.RunNextInstruction(); cpu.cycles = 0
  // STA $24
  cpu.RunNextInstruction(); cpu.cycles = 0

  // LDX $24
  cpu.RunNextInstruction();
  if cpu.cycles != 3 {
    log.Printf("Expecting 3 CPU cycles but got %d", cpu.cycles)
    t.Fail()
  }
  if cpu.Y() != 42 {
    log.Printf("Expecting register X to be 42 but got %d", cpu.X())
    t.Fail()
  }
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
    // Zero Page
    0xA9, 42,   // LDA #42
    0x85, 0x24, // STA $24

    // Zero Page,X
    0xA9, 24,   // LDA #24
    0xA2, 3,    // LDX #3
    0x95, 0x42, // STA $42,X
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

  cpu.RunNextInstruction(); cpu.cycles = 0
  cpu.RunNextInstruction(); cpu.cycles = 0
  cpu.RunNextInstruction()
  if cpu.cycles != 4 {
    log.Printf("Expecting CPU cycles to be %d, but got %d", 4, cpu.cycles)
    t.Fail()
  }
  if cpu.memory.GetUint8At(0x42 + 3) != 24 { t.Fail() }
  cpu.cycles = 0
}

func xTestStx(t *testing.T) {
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