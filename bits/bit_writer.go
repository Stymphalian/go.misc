package bits

import (
	"io"
)

type BitWriter interface {
	// Write a single bit out to the byte stream.
	// Args:
	//   is_one: 1 for writing a '1' else 0 for writing a '0'
	// Return:
	//   error: Return an error if something failed to write
	Write(is_one int) error

	// BitWriter can only write at a byte level granularity. We flush the
	// last byte to the stream with the remaining valid 'bits'. The last
	// byte will be padded with 0's
	// Return:
	//   error: Error if anything goes wrong.
	Flush() error

	// Returns the number of bits that have been written.
	// Return:
	//   int: number of bits.
	NumBitsWritten() int
}

type BitsWriter struct {
	// A counter for the number of bits which have been written.
	NumBits int

	// The byte stream to write to. We can only write at a byte level
	// granularity so the last byte may potentially be zero padded.
	OutputStream io.ByteWriter

	// The index into the accumulationByte in which to write the next bit.
	// Must be a value between 7 -> 0
	bitPos uint

	// A scratch byte where the bits are set. Once all 8 bits have been filled,
	// it is then wrote out into the byte stream
	accumulationByte byte
}

// Create a new BitsWriter.
// Args:
//   stream: The byte stream in which to write the bits to. Note the bit writer
//     can only write at a byte level granularity.
//     For Example: 
//       If we only have the bits '101' to write in the last byte.
//       The full byte '1010 0000' will still be output to the stream.
// Return:
//   *BitsWriter: the bit writer object
func NewBitsWriter(stream io.ByteWriter) *BitsWriter {
	return &BitsWriter{0, stream, uint(7), 0x00}
}

// Write a single bit to the stream.
// Bits are written to the byte from left -> right.
// Args:
//   is_one: 1 if we want to write a '1', else 0 if we want to write '0'
// Return:
//   error: nil if succeeds, else an error
func (self *BitsWriter) Write(is_one int) error {
	self.NumBits += 1

	if is_one == 1 {
		self.accumulationByte |= (1 << self.bitPos)
	}

	self.bitPos -= 1
	if int(self.bitPos) < 0 {
		// we have finsihed filling the scratch byte. Dump it out to the stream.
		err := self.OutputStream.WriteByte(self.accumulationByte)
		if err != nil {
			return err
		}

		// reset the bit position
		self.bitPos = uint(7)
		// Note that the byte is reset to all 0's
		self.accumulationByte = 0x00
	}

	return nil
}

// Return the number of bits written. This count does not include any
// padded 0's of the LAST byte which may have been outputted due to the
// byte level granularity of the BitWriter
// Return:
//   int: The number of bits written.
func (self *BitsWriter) NumBitsWritten() int {
	return self.NumBits
}

// Flush the remaining byte to the output stream.
// Return:
//   error: If anything goes wrong during the write.
func (self *BitsWriter) Flush() error {
	if self.bitPos != uint(7) {
		err := self.OutputStream.WriteByte(self.accumulationByte)
		if err != nil {
			return err
		}
	}
	return nil
}
