package main

const (
  // The tags
  C byte = 1 << iota
  Z
  I
  D
  B
  _
  V
  N
)

type CPU struct {
  // Registers
  A, X, Y, P byte
}

func CPUNew() CPU {
  return CPU{ A:0, X:0, Y:0, P: 0 }
}

func (c* CPU) C() bool {
  return (c.P & C) != 0
}

func (c* CPU) setC(status bool) {
  if status {
    c.P = c.P | C
  } else {
    c.P = c.P & ^C
  }
}

func (c* CPU) Z() bool {
  return (c.P & Z) != 0
}

func (c* CPU) setZ(status bool) {
  if status {
    c.P = c.P | Z
  } else {
    c.P = c.P & ^Z
  }
}