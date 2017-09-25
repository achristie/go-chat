package msg_test

import "testing"

// TestEncode test that the encoding of a message works.
func TestEncode(t *testing.T) {
	m := msg.MSG{
		Name: "0123456789",
		Data: "hello",
	}

	data := msg.Encode(m)

	if len(data) != 17 {
		t.Fatalf("X - not correct number of bytes, exp[17] got[%d]\n", len(data))
	}
	t.Log("correct number of bytes")
}
