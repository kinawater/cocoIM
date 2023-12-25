package websocket

import (
	"cocoIM"
	"github.com/gobwas/ws"
	"net"
)

// 把ws frame封装成应用层的私有ws frame
// ws Frame结构
//  0                   1                   2                   3
//  0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
// +-+-+-+-+-------+-+-------------+-------------------------------+
// |F|R|R|R| opcode|M| Payload len |    Extended payload length    |
// |I|S|S|S|  (4)  |A|     (7)     |             (16/64)           |
// |N|V|V|V|       |S|             |   (if payload len==126/127)   |
// | |1|2|3|       |K|             |                               |
// +-+-+-+-+-------+-+-------------+ - - - - - - - - - - - - - - - +
// |     Extended payload length continued, if payload len == 127  |
// + - - - - - - - - - - - - - - - +-------------------------------+
// |                               |Masking-key, if MASK set to 1  |
// +-------------------------------+-------------------------------+
// | Masking-key (continued)       |          Payload Data         |
// +-------------------------------- - - - - - - - - - - - - - - - +
// :                     Payload Data continued ...                :
// + - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - +
// |                     Payload Data continued ...                |
// +---------------------------------------------------------------+
//
//其中，各字段含义如下：
//
//● FIN：1位，表示是否是最后一个数据帧。
//● RSV1, RSV2, RSV3：各占1位，保留位。
//● Opcode：4位，表示数据类型，如文本（0x1）、二进制（0x2）等。
//● MASK：1位，表示是否使用掩码。
//● Payload length：7位或7+16位或7+64位，表示负载数据的长度。
//● Masking-key：4字节，如果MASK位被设置为1，那么就存在4字节的掩码值。
//● Payload data：负载数据，长度等于Payload length指定的长度。
//需要注意的是，如果Payload length为126，则后续两个字节表示的是负载数据的实际长度。
//如果Payload length为127，则后续八个字节表示的是负载数据的实际长度。
//如果MASK位被设置为1，则接下来的4个字节是用来对负载数据进行掩码的。

type Frame struct {
	WsFrame ws.Frame
}

func (f *Frame) SetOpCode(code cocoIM.OpCode) {
	f.WsFrame.Header.OpCode = ws.OpCode(code)
}

func (f *Frame) GetOpCode() cocoIM.OpCode {
	return cocoIM.OpCode(f.WsFrame.Header.OpCode)
}

func (f *Frame) SetPayload(payloadBytes []byte) {
	f.WsFrame.Payload = payloadBytes
}

func (f *Frame) GetPayload() []byte {
	// 有掩码
	if f.WsFrame.Header.Masked {
		// 解码
		ws.Cipher(f.WsFrame.Payload, f.WsFrame.Header.Mask, 0)
	}
	f.WsFrame.Header.Masked = false
	return f.WsFrame.Payload
}

type WsConn struct {
	net.Conn
}
