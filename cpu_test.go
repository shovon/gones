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

func TestAdc(t *testing.T) {
  
}

func TestLda(t *testing.T) {
  instructions := ConvertSimpleInstructions([]byte{
    0xA9, 3,
    0xA9, 128,
  })
  cpu := CPUNew()
  cpu.SetInstructions(instructions)
  cpu.MovePCToResetVector()
  if cpu.pc != 0x8000 {
    log.Printf("Expecting program counter to point to location 0x8000, but it's pointing at %d", cpu.pc)
    t.Fail()
  }

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
}

func TestSta(t *testing.T) {

}
