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

func TestZ(t *testing.T) {
  cpu := CPUNew()
  if (cpu.Z()) {
    t.Fail()
  }

  cpu.setZ(true)

  if (!cpu.Z()) {
    t.Fail()
  }

  cpu.setZ(false)

  if (cpu.Z()) {
    t.Fail()
  }
}