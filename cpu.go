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

func (c* CPU) status(flag byte) bool {
  return (c.P & flag) != 0
}

func (c* CPU) C() bool {
  return c.status(C)
}

func (c* CPU) SetC(status bool) {
  if status {
    c.P = c.P | C
  } else {
    c.P = c.P & ^C
  }
}

func (c* CPU) Z() bool {
  return c.status(Z)
}

func (c* CPU) SetZ(status bool) {
  if status {
    c.P = c.P | Z
  } else {
    c.P = c.P & ^Z
  }
}