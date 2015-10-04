// bits provides interfaces/objects for reading/writing bits to a stream.
package bits

import (
	"errors"
	"io"
)

type BitReader interface {
	// Read a bit. Return 1  or 0
	// Return:
	//   int: 1 or 0
	//   error: nil if successful. Not nil if there was an error, or this is
	//     the last bit of the stream.
	Read() (int, error)
}

type BitsReader struct {
	// The byte stream to read from
	Stream io.ByteReader

	// The number of bits to skip in the LAST byte.
	SkipBits uint

	// A look-ahead byte.
	nextByte byte
	// A flag to denote if the 'bufferedByte' is the LAST byte of the stream.
	hasNextByte bool
	// The position the read in currently indexed into the buffered bytes
	// Must be a value between 7 -> 0
	bitPos uint
	// The current byte being read from.
	bufferedByte byte
}

// Create a new BitsReader from the given ByteReader
// Args:
//   r : Reads each byte and then return the bits individually
//   num_trash_bits: num bits of the LAST byte of the stream to skip.
//	   For Example.
//       Given num_trash_bits = 3
//			 and a last byte = 1110 1000
//			 The BitReader will return 11101.
//       The last set of 000 is not returned.
// Return:
//   *BitReader object
func NewBitsReader(r io.ByteReader, num_trash_bits uint) *BitsReader {
	if num_trash_bits >= 8 {
		panic("num_trash_bits must be < 8")
	}
	return &BitsReader{r, num_trash_bits, 0x00, false, uint(7), 0x00}
}

// Read a single bit from the byte stream
// Returns:
//   int: 1 or 0 for the bit value
//   error: nil if there is not error. The read will return an error if there
//	   are no more bits in the stream to read.
func (self *BitsReader) Read() (int, error) {
	if self.bitPos == uint(7) {
		// The bit position is at the beginning, we must first fetch a new byte
		// to start reading from.
		var b byte
		var err error

		if self.hasNextByte {
			// The BitsReader stores one look-ahead byte. If we have already read
			// this byte from the byte stream then just use it.
			b = self.nextByte
		} else {
			// This is the first time retrieving a byte from the stream.
			b, err = self.Stream.ReadByte()
			if err != nil {
				return 0, err
			}
		}

		// Try to get the next look-ahead byte.
		// If there was an error,then we know there are no more bytes left in the
		// stream, and that the current 'bufferedByte' is the LAST byte.
		nextByte, err := self.Stream.ReadByte()
		self.hasNextByte = (err == nil)
		self.nextByte = nextByte

		// save the byte for reading the bits off of.
		self.bufferedByte = b
	}

	if !self.hasNextByte && self.bitPos < self.SkipBits {
		// This is the last byte, and we have read all the 'valid' bits of the byte
		// Return with an error singalling that the reading is finished.
		return 0, errors.New("Finished reading")
	}

	// get the bit value, and advance the bitpos pointer
	bit := self.bufferedByte & (1 << self.bitPos)
	self.bitPos -= 1

	// reset the position to the start of the next byte
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
