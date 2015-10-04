package bits

import (
	"bytes"
	"io"
	"testing"
)

func TestBitsReader(t *testing.T) {
	cases := []struct {
		data           []byte
		num_trash_bits uint
		want           string
	}{
		{[]byte{0xA0}, 0, "10100000"},
		{[]byte{0x80}, 0, "10000000"},
		{[]byte{0xAA}, 0, "10101010"},
		{[]byte{0xFF}, 0, "11111111"},
		{[]byte{0x00}, 0, "00000000"},
		{[]byte{0x2C}, 0, "00101100"},
		{[]byte{0x2C, 0xA0}, 0, "0010110010100000"},
		{[]byte{0x2C, 0xA8}, 3, "0010110010101"},
		{[]byte{0x2C, 0xA8}, 7, "001011001"},
		{[]byte{0x2C, 0xAE}, 1, "001011001010111"},
	}

	for i, c := range cases {
		var get bytes.Buffer
		stream := bytes.NewBuffer(c.data)
		r := NewBitsReader(io.ByteReader(stream), c.num_trash_bits)

		for {
			bit, err := r.Read()

			// we have reached the end of the bit stream
			if err != nil {
				break
			}

			// write out the bit to a string so that we can compare later
			if bit == 1 {
				get.WriteString("1")
			} else {
				get.WriteString("0")
			}
		}

		// compare the string we got with the one we want.
		get_string := get.String()
		if get_string != c.want {
			t.Errorf("[Test %d] get != want\nget :%v\nwant:%v\n", i, get_string,
				c.want)
		}
	}

}
