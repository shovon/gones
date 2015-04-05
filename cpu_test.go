package main

import "testing"

func TestC(t *testing.T) {
  cpu := CPUNew()
  if (cpu.C()) {
    t.Fail()
  }

  cpu.SetC(true)

  if (!cpu.C()) {
    t.Fail()
  }

  cpu.SetC(false)

  if (cpu.C()) {
    t.Fail()
  }
}

func TestZ(t *testing.T) {
  cpu := CPUNew()
  if (cpu.Z()) {
    t.Fail()
  }

  cpu.SetZ(true)

  if (!cpu.Z()) {
    t.Fail()
  }

  cpu.SetZ(false)

  if (cpu.Z()) {
    t.Fail()
  }
}