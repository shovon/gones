package main

import "testing"

func TestC(t *testing.T) {
  cpu := CPUNew()
  if (cpu.C()) {
    t.Fail()
  }
}