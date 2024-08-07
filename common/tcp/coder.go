package tcp

import (
	"bytes"
	"encoding/binary"
)

type DataPacket struct {
	Len  uint32
	Data []byte
}

// Marshal 对数据包进行序列化，主要是两部分，前4个字节是整个消息包的长度，从第五个字节开始是消息包
// ----------------------------------------------------------------------
// | 1 | 2 | 3 | 4 |  ... ...                                            |
// _______________________________________________________________________
// |    数据长度标识  |    数据                                             |
// ----------------------------------------------------------------------
func (d *DataPacket) Marshal() ([]byte, error) {
	dataBuffer := bytes.NewBuffer([]byte{})
	err := binary.Write(dataBuffer, binary.BigEndian, d.Len)
	if err != nil {
		return nil, err
	}
	return append(dataBuffer.Bytes(), d.Data...), nil
}
