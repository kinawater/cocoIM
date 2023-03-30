package wsCommProto

import (
	"fmt"
	"testing"
)

func TestCommProto_MarshalMarshal(t *testing.T) {
	message := CommProto{
		Command:   0x1234,
		MsgLength: 8,
		Message:   []byte{0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7, 0x8},
	}
	marshal, err := message.Marshal()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%x\n", marshal)

	packetMsg := CommProto{}
	err = packetMsg.UnMarshal(marshal)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v\n", packetMsg)

	message2 := CommProto{
		Command:   101,
		MsgLength: 0,
		Message:   []byte{},
	}
	marshal2, _ := message2.Marshal()
	fmt.Printf("%x\n", marshal2)
	packetMsg2 := CommProto{}
	_ = packetMsg2.UnMarshal(marshal2)
	fmt.Printf("%v\n", packetMsg2)
}
