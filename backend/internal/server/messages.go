package server

import "encoding/json"

const (
	MessageTypeLoadDocument    MessageType = "load-document"
	MessageTypeUpdateDocument  MessageType = "update-document"
	MessageTypeSaveDocument    MessageType = "save-document"
	MessageTypeUpdateUserCount MessageType = "update-user-count"
)

type MessageType string

type Message struct {
	Type    MessageType     `json:"type"`
	Payload json.RawMessage `json:"payload"`
}
