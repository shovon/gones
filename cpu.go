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

func (c* CPU) Status(flag byte) bool {
  return (c.P & flag) != 0
}

func (c* CPU) SetStatus(flag byte, status bool) {
  if status {
    c.P = c.P | flag
  } else {
    c.P = c.P & ^flag
  }
}
