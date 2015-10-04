package bits

import (
	"bytes"
	"io"
	"testing"
)

func compareByteSlices(a []byte, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i += 1 {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestBitsWriter(t *testing.T) {
	cases := []struct {
		bits string
		want []byte
	}{
		{"01110000", []byte{0x70}},
		{"10100101", []byte{0xA5}},
		{"0111000010100101", []byte{0x70, 0xA5}},
		{"01110000101", []byte{0x70, 0xA0}},
		{"", []byte{}},
		{"1", []byte{0x80}},
		{"0", []byte{0x00}},
		{"1110", []byte{0xE0}},
	}

	for i, c := range cases {
		// create a new BitsWrite, we will write to a bytes.Buffer
		var output bytes.Buffer
		w := NewBitsWriter(io.ByteWriter(&output))

		// make sure we start with 0 bits written to the stream
		if w.NumBitsWritten() != 0 {
			t.Errorf("[Test %d] writer should be default to 0 bits", i)
		}

		// iterate throught the case string and output a 0 when we see "0"
		// and a 1 when we see "1"
		var err error
		for j := 0; j < len(c.bits); j += 1 {
			if c.bits[j] == "0"[0] {
				err = w.Write(0)
			} else {
				err = w.Write(1)
			}

			// check for eany write errors
			if err != nil {
				t.Errorf("[Test %d] failed to write bit", i)
			}
		}

		// Flush any remaining bits into the output stream
		err = w.Flush()
		if err != nil {
			t.Errorf("[Test %d] failed to write bit", i)
		}

		// Check that the proper number of bits were written
		if len(c.bits) != w.NumBitsWritten() {
			t.Errorf("[Test %d] not enough bits written to stream\n%v outof %v bits",
				i, w.NumBitsWritten(), len(c.bits))
		}

		// compare the expected vs recevied []byte slices
		get := output.Bytes()
		if !compareByteSlices(get, c.want) {
			t.Errorf("[Test %d] Bytes slices are different!\ngot :%v\nwant:%v\n", i, get,
				c.want)
		}
	}

}
