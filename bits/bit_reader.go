package bits

import (
	"errors"
	"io"
)

type BitReader interface {
	// Read a bit. Return 1  or 0
	Read() (int, error)
}

type BitsReader struct {
	Stream   io.ByteReader
	SkipBits uint

	nextByte     byte
	hasNextByte  bool
	bitPos       uint
	bufferedByte byte
}

// Create a new BitsReader from the given ByteReader
// Args:
//   r : Reads each byte and then return ths bits individually
//   num_trash_bits: num bits of the last byte which are invalid
// Return:
//   *BitReader object
func NewBitsReader(r io.ByteReader, num_trash_bits uint) *BitsReader {
	if num_trash_bits >= 8 {
		panic("num_trash_bits must be < 8")
	}
	return &BitsReader{r, num_trash_bits, 0x00, false, uint(7), 0x00}
}

func (self *BitsReader) Read() (int, error) {
	// This it the beginning of the read
	if self.bitPos == uint(7) {
		var b byte
		var err error

		if self.hasNextByte {
			b = self.nextByte

			nextByte, err := self.Stream.ReadByte()
			self.hasNextByte = (err == nil)
			self.nextByte = nextByte
		} else {
			b, err = self.Stream.ReadByte()
			if err != nil {
				return 0, err
			}

			nextByte, err := self.Stream.ReadByte()
			self.hasNextByte = (err == nil)
			self.nextByte = nextByte
		}

		self.bufferedByte = b
	}

	// get the bit value, and advance the bitpos pointer
	if !self.hasNextByte && self.bitPos < self.SkipBits {
		return 0, errors.New("Finished reading")
	}

	bit := self.bufferedByte & (1 << self.bitPos)
	self.bitPos -= 1
	if int(self.bitPos) < 0 {
		self.bitPos = 7
	}

	// return the bit value
	if bit != 0 {
		return 1, nil
	} else {
		return 0, nil
	}
}
