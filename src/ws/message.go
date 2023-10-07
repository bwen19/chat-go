package ws

import (
	"encoding/json"
	"errors"
)

type Message struct {
	Type  string
	Event interface{}
}

type SendMessageEvent struct {
	RoomID  int64
	Content string
	Kind    string
}

func (s *Server) handleMessage(jsonMessage []byte) error {
	var event json.RawMessage
	msg := Message{Event: &event}
	if err := json.Unmarshal(jsonMessage, &msg); err != nil {
		return err
	}

	switch msg.Type {
	case "send-message":
		var e SendMessageEvent
		if err := json.Unmarshal(event, &e); err != nil {
			return err
		}

	default:
		return errors.New("unknown message type")
	}
	return nil
}
