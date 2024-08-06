package tcp

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"time"
)

func ReadData(conn *net.TCPConn) ([]byte, error) {
	var dataLen uint32
	dataLenBuf := make([]byte, 4)
	if err := readFixedData(conn, dataLenBuf); err != nil {
		return nil, err
	}
	// 构造一个缓冲区
	buffer := bytes.NewBuffer(dataLenBuf)
	if err := binary.Read(buffer, binary.BigEndian, &dataLen); err != nil {
		return nil, err
	}
	// 读取长度标识符后面的数据
	if dataLen <= 0 {
		return nil, fmt.Errorf("收到的消息长度为%d，需要检查", dataLen)
	}
	var dataBuf = make([]byte, dataLen)
	if err := readFixedData(conn, dataBuf); err != nil {
		return nil, fmt.Errorf("read headlen error:%s", err.Error())
	}
	return dataBuf, nil

}

func readFixedData(conn *net.TCPConn, buf []byte) error {
	// 设置超时时间
	_ = (*conn).SetReadDeadline(time.Now().Add(time.Duration(120) * time.Second))

	var pos = 0
	// 提前设定好buf长度，buf的长度就是要读取的长度
	var totalSize = len(buf)
	for {
		readDataLen, err := (*conn).Read(buf[pos:])
		if err != nil {
			return err
		}
		pos = pos + readDataLen
		// 已经读到足够的数据
		if pos == totalSize {
			break
		}
	}
	return nil

}
