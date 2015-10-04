package bits

import (
	"io"
	//"log"
)

type BitWriter interface {
	Write(is_one int) error
	Flush() error
	NumBitsWritten() int
}

type BitsWriter struct {
	NumBits      int
	OutputStream io.ByteWriter

	bitPos           uint
	accumulationByte byte
}

func NewBitsWriter(stream io.ByteWriter) *BitsWriter {
	return &BitsWriter{0, stream, uint(7), 0x00}
}

func (self *BitsWriter) Write(is_one int) error {
	self.NumBits += 1

	if is_one == 1 {
		self.accumulationByte |= (1 << self.bitPos)
	}

	self.bitPos -= 1
	if int(self.bitPos) < 0 {
		err := self.OutputStream.WriteByte(self.accumulationByte)
		if err != nil {
			return err
		}

		// reset the bit position
		self.bitPos = uint(7)
		self.accumulationByte = 0x00
	}

	return nil
}

func (self *BitsWriter) NumBitsWritten() int {
	return self.NumBits
}

func (self *BitsWriter) Flush() error {
	if self.bitPos != uint(7) {
		err := self.OutputStream.WriteByte(self.accumulationByte)
		if err != nil {
			return err
		}
	}
	return nil
}
