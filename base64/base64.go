// base64 provides Encode and Decode functions for base64 encoding.
package base64

import (
	"bytes"
)

const (
	codes = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/="
)

// Helper function for finding the index of a string byte in the 'codes' const
// Args:
//   index: The 'letter' in the codes in which to retrieve the index for
// Return:
//   byte: The int position of the index value returned as a byte.
func indexOf(index byte) byte {
	if index >= byte("A"[0]) && index <= byte("Z"[0]) {
		return index - byte("A"[0])
	} else if index >= byte("a"[0]) && index <= byte("z"[0]) {
		return index - byte("a"[0]) + 26
	} else if index >= byte("0"[0]) && index <= byte("9"[0]) {
		return index - byte("0"[0]) + 52
	}
	if index == byte("+"[0]) {
		return byte(len(codes) - 3)
	}
	if index == byte("/"[0]) {
		return byte(len(codes) - 2)
	}
	return byte(len(codes) - 1)
}

// Function to encode a blob of binary data into a base64 encoded string
// Args:
//   data: The binary data in which to encode
// Return:
//   string: The base64 encoded string.
func Encode(data []byte) string {
	var buffer bytes.Buffer
	var b byte
	var remainder byte
	state := 1
	data_len := len(data)

	// |  1  |   2   |  3  |
	// 6 - 2 - 4 - 4 - 2 - 6
	for i := 0; i < data_len; i++ {
		if state == 1 {
			//6-2
			b = (data[i] & 0xFC) >> 2
			buffer.WriteString(string(codes[b]))
			remainder = (data[i] & 0x03) << 4
			state = 2

		} else if state == 2 {
			//4-4
			b = remainder
			b |= (data[i] & 0xF0) >> 4
			buffer.WriteString(string(codes[b]))
			remainder = (data[i] & 0x0F) << 2
			state = 3

		} else if state == 3 {
			//2-6
			b = remainder
			b |= ((data[i] & 0xC0) >> 6)
			buffer.WriteString(string(codes[b]))

			b = data[i] & 0x3F
			buffer.WriteString(string(codes[b]))
			state = 1
		}
	}

	// Handle the end byte cases.
	if state == 2 {
		// 1 byte of data, 6 bits + 2 bits leftover
		buffer.WriteString(string(codes[remainder]))
		buffer.WriteString("==")
	} else if state == 3 {
		// 2 bytes of data, 6 + 2 + 4, 4 bits leftover
		buffer.WriteString(string(codes[remainder]))
		buffer.WriteString("=")
	}

	// return the string.
	return buffer.String()
}

// Function to decode a base64 encoded string into a blob of binary data.
// Args:
//   data: The base64 encoded string to decode.
// Return:
//   []byte: The decoded binary
func Decode(data string) []byte {
	if len(data)%4 != 0 {
		panic("Not a base64 encoded string")
	}
	if len(data) == 0 {
		return make([]byte, 0)
	}

	var buffer bytes.Buffer
	i := 0
	data_len := len(data)

	for i = 0; i < data_len-4; i += 4 {
		b1 := indexOf(data[i])
		b2 := indexOf(data[i+1])
		b3 := indexOf(data[i+2])
		b4 := indexOf(data[i+3])

		// |  1  |   2   |  3  |
		// 6 - 2 - 4 - 4 - 2 -6
		buffer.WriteByte((b1 << 2) | ((b2 & 0x30) >> 4))
		buffer.WriteByte((b2 << 4) | ((b3 & 0x3C) >> 2))
		buffer.WriteByte((b3 << 6) | (b4 & 0x3F))
	}

	// process the last 4 characters
	b1 := indexOf(data[i])
	b2 := indexOf(data[i+1])
	b3 := indexOf(data[i+2])
	b4 := indexOf(data[i+3])
	if data[data_len-1] == byte("="[0]) {
		if data[data_len-2] == byte("="[0]) {
			// Case with "==" at the end
			buffer.WriteByte((b1 << 2) | ((b2 & 0x30) >> 4))
		} else {
			// Case with "=" at the end.
			// 2 bytes of data.
			buffer.WriteByte((b1 << 2) | ((b2 & 0x30) >> 4))
			buffer.WriteByte((b2 << 4) | ((b3 & 0x3C) >> 2))
		}
	} else {
		buffer.WriteByte((b1 << 2) | ((b2 & 0x30) >> 4))
		buffer.WriteByte((b2 << 4) | ((b3 & 0x3C) >> 2))
		buffer.WriteByte((b3 << 6) | (b4 & 0x3F))
	}

	return buffer.Bytes()
}
