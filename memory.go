package main

const MEMORY_SIZE uint16 = 1024*64 - 1

type Memory [MEMORY_SIZE]byte

// Gets two contiguous bytes at the specified memory location, interpreting them
// as little endian 16-bit integers.
func (m *Memory) GetUint16LEAt(location uint16) uint16 {
  return (uint16(m[location+1]) << 8) | uint16(m[location])
}

// Utility function to grab a 16-bit value
func GetUint16LEAt(buffer []byte, location uint16) uint16 {
  lsb := buffer[location]
  msb := buffer[location+1]
  value := (uint16(msb) << 8) | uint16(lsb)
  return value
}

// Gets the 8-bit integer at the specified memory location.
func (m *Memory) GetUint8At(location uint16) byte {
  return m[location]
}

// Sets the specified memory location with an 8-bit integer.
func (m *Memory) SetUint8At(location uint16, value byte) {
  m[location] = value
}

// Sets the memory with the instructions.
func (m *Memory) SetInstructions(instructions []byte) {
  copy(m[0x8000:], instructions)
}
