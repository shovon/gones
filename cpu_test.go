package main

import "testing"

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

func TestAdc(t *testing.T) {

}