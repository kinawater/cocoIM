package wsCommProto

import (
	"bytes"
	"encoding/binary"
	"errors"
)

const protoHeaderMinSize = 6 //协议最小长度

type CommProto struct {
	Command   uint16 // 2 bytes  指令
	MsgLength uint32 // 4 bytes 消息长度
	Message   []byte //真正的消息载体
}

func (p *CommProto) Marshal() ([]byte, error) {
	b := new(bytes.Buffer)
	err := binary.Write(b, binary.BigEndian, p.Command)
	if err != nil {
		return nil, err
	}
	err = binary.Write(b, binary.BigEndian, p.MsgLength)
	if err != nil {
		return nil, err
	}
	_, err = b.Write(p.Message)
	if err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}
func (p *CommProto) UnMarshal(message []byte) error {
	if len(message) < protoHeaderMinSize {
		return errors.New("message的长度还不足一个空协议的长度，协议不完整，请检查")
	}
	i := 0
	command := binary.BigEndian.Uint16(message[i : i+2])
	i += 2
	payloadLen := binary.BigEndian.Uint32(message[i : i+4])
	if cap(p.Message) < int(payloadLen) {
		p.Message = make([]byte, payloadLen)
	} else {
		p.Message = p.Message[:payloadLen]
	}
	p.Command = command
	p.MsgLength = payloadLen
	copy(p.Message, message[protoHeaderMinSize:protoHeaderMinSize+payloadLen])
	return nil
}
