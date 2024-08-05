package gateway

import "net"

type connection struct {
	fd   int
	conn *net.TCPConn
}

func (c *connection) Close() error {
	err := c.conn.Close()
	if err != nil {
		return err
	}
	return nil
}

// 获取地址
func (c *connection) RemoteAddr() string {
	return c.conn.RemoteAddr().String()
}
