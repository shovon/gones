package main

const MEMORY_SIZE uint16 = 1024*64 - 1

type Memory [MEMORY_SIZE]byte

// Gets the least significant byte in the location at 
func (m *Memory) GetUint16LEAt(location uint16) uint16 {
  return (uint16(m[location+1]) << 8) & uint16(m[location])
}

func (m *Memory) GetInt8At(location uint16) byte {
  return m[location]
}

func (m *Memory) SetInstructions(instructions []byte) {
  copy(m[0x8000:0xC000], instructions)
}
