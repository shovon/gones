package main

import "testing"

func testStatus(t *testing.T, flag byte) {
  cpu := CPUNew()
  if (cpu.Status(flag)) {
    t.Fail()
  }

  cpu.SetStatus(flag, true)

  if (!cpu.Status(flag)) {
    t.Fail()
  }

  cpu.SetStatus(flag, false)

  if (cpu.Status(flag)) {
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
