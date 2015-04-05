package main

import "testing"

func TestC(t *testing.T) {
  cpu := CPUNew()
  if (cpu.C()) {
    t.Fail()
  }

  cpu.setC(true)

  if (!cpu.C()) {
    t.Fail()
  }

  cpu.setC(false)

  if (cpu.C()) {
    t.Fail()
  }
}
