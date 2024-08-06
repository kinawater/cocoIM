package tcp

import (
	"bytes"
	"encoding/binary"
)

type DataPacket struct {
	Len  uint32
	Data []byte
}

func (d *DataPacket) Marshal() ([]byte, error) {
	dataBuffer := bytes.NewBuffer([]byte{})
	err := binary.Write(dataBuffer, binary.BigEndian, d.Len)
	if err != nil {
		return nil, err
	}
	return append(dataBuffer.Bytes(), d.Data...), nil
}
