package source

import (
	"cocoIM/common/discovery"
	"fmt"
)

type EventType string

type Event struct {
	Type         EventType
	IP           string
	Port         string
	ConnectNum   float64
	MessageBytes float64
}

const (
	AddNodeEvent EventType = "addNode"
	DelNodeEvent EventType = "delNode"
)

// “事件”channel
var eventChan chan *Event

func EventChan() <-chan *Event { return eventChan }

func NewEvent(ed *discovery.EndpointInfo) *Event {
	if ed == nil || ed.MetaData == nil {
		return nil
	}
	var connNum, msgBytes float64

	if data, ok := ed.MetaData["connect_num"]; ok {
		connNum = data.(float64)
	}

	if data, ok := ed.MetaData["message_bytes"]; ok {
		msgBytes = data.(float64)
	}
	return &Event{
		Type:         AddNodeEvent,
		IP:           ed.IP,
		Port:         ed.Port,
		ConnectNum:   connNum,
		MessageBytes: msgBytes,
	}
}
func (nd *Event) Key() string {
	return fmt.Sprintf("%s:%s", nd.IP, nd.Port)
}
