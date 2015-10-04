package base64

import (
	"reflect"
	"testing"
)

var kCases = []struct {
	raw     []byte
	encoded string
}{
	{[]byte(""), ""},
	{[]byte("a"), "YQ=="},
	{[]byte("aa"), "YWE="},
	{[]byte("abc"), "YWJj"},
	{[]byte{0xFC}, "/A=="},
	{[]byte("Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum"),
		"TG9yZW0gaXBzdW0gZG9sb3Igc2l0IGFtZXQsIGNvbnNlY3RldHVyIGFkaXBpc2NpbmcgZWxpdCwgc2VkIGRvIGVpdXNtb2QgdGVtcG9yIGluY2lkaWR1bnQgdXQgbGFib3JlIGV0IGRvbG9yZSBtYWduYSBhbGlxdWEuIFV0IGVuaW0gYWQgbWluaW0gdmVuaWFtLCBxdWlzIG5vc3RydWQgZXhlcmNpdGF0aW9uIHVsbGFtY28gbGFib3JpcyBuaXNpIHV0IGFsaXF1aXAgZXggZWEgY29tbW9kbyBjb25zZXF1YXQuIER1aXMgYXV0ZSBpcnVyZSBkb2xvciBpbiByZXByZWhlbmRlcml0IGluIHZvbHVwdGF0ZSB2ZWxpdCBlc3NlIGNpbGx1bSBkb2xvcmUgZXUgZnVnaWF0IG51bGxhIHBhcmlhdHVyLiBFeGNlcHRldXIgc2ludCBvY2NhZWNhdCBjdXBpZGF0YXQgbm9uIHByb2lkZW50LCBzdW50IGluIGN1bHBhIHF1aSBvZmZpY2lhIGRlc2VydW50IG1vbGxpdCBhbmltIGlkIGVzdCBsYWJvcnVt"},
}

func TestEncode(t *testing.T) {
	for i, c := range kCases {
		get := Encode(c.raw)
		if get != c.encoded {
			t.Errorf("Failed at index %v: %v", i, get)
		}
	}
}

func TestDecode(t *testing.T) {
	for i, c := range kCases {
		get := Decode(c.encoded)
		if !reflect.DeepEqual(get, c.raw) {
			t.Errorf("Failed at index %v", i)
		}
	}
}

func TestIndexOf(t *testing.T) {
	cases := []struct {
		letter byte
		want   byte
	}{
		{byte("A"[0]), 0},
		{byte("B"[0]), 1},
		{byte("Z"[0]), 25},
		{byte("a"[0]), 26},
		{byte("b"[0]), 27},
		{byte("z"[0]), 51},
		{byte("0"[0]), 52},
		{byte("1"[0]), 53},
		{byte("9"[0]), 61},
		{byte("+"[0]), 62},
		{byte("/"[0]), 63},
		{byte("="[0]), 64},
	}

	for i, c := range cases {
		get := indexOf(c.letter)
		if get != c.want {
			t.Errorf("Failed at index %v", i)
		}
	}
}
