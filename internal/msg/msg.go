package msg

import (
	"encoding/binary"
	"io"

	"github.com/pkg/errors"
)

const hdrLength = 12

// MSG defines the message protocol data.
type MSG struct {
	Name string
	Data string
}

// Read waits on the network to receive a chat message.
func Read(r io.Reader) ([]byte, int, error) {

	// Read the first header length of bytes.
	buf := make([]byte, hdrLength)
	if _, err := io.ReadFull(r, buf); err != nil {
		errors.Wrap(err, "ReadFull header")
		return nil, 0, err
	}

	// Get the length for the remaining bytes.
	length := int(binary.BigEndian.Uint16(buf[10:12])) + hdrLength

	// Copy the header bytes into the final slice.
	data := make([]byte, length)
	copy(data, buf)

	// read the remaning bytes.
	if _, err := io.ReadFull(r, data[hdrLength:]); err != nil {
		errors.Wrap(err, "ReadFull data")
		return nil, 0, err
	}

	return data, length, nil
}

// Decode will take the bytes
func Decode(data []byte) MSG {
	msg := MSG{
		Name: string(data[:10]),
		Data: string(data[12:]),
	}
}

// Encode will take a message and produce byte slice.
func Encode(msg MSG) []byte {
	data := make([]byte, hdrLength*len(msg.Data))

	copy(data, msg.Name[:10])
	binary.BigEndian.PutUint16(data[10:12], uint16(len(msg.Data)))
	copy(data[12:], msg.Data)

	return data
}
