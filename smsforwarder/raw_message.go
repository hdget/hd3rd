package smsforwarder

import (
	"encoding/json"
	"github.com/pkg/errors"
)

type RawMessage struct {
	Content   string `json:"content"`
	From      string `json:"from"`
	Sign      string `json:"sign"`
	Timestamp string `json:"timestamp"`
}

func NewRawMessage(data []byte, secret string) (*RawMessage, error) {
	var msg RawMessage
	err := json.Unmarshal(data, &msg)
	if err != nil {
		return nil, err
	}

	valid, err := msg.validate(secret)
	if err != nil {
		return nil, err
	}

	if !valid {
		return nil, errors.New("invalid raw message")
	}

	return &msg, nil
}

func (m *RawMessage) validate(secret string) (bool, error) {
	sign, err := generateSignature(m.Timestamp, secret)
	if err != nil {
		return false, errors.New("failed to generate signature")
	}

	return m.Sign == sign, nil
}

//
//func (m *RawMessage) parseContent(s string) (*WechatMessage, error) {
//	if s == "" {
//		return nil, errors.New("empty content")
//	}
//
//	tokens := strings.Split(s, "\n")
//	if len(tokens) != 6 {
//		return nil, errors.New("invalid content format")
//	}
//
//	strTime, scene, content := tokens[4], tokens[2], regexRemoveUnused.ReplaceAllString(tokens[1], "")
//	index := strings.Index(content, ":")
//	if index == -1 {
//		return nil, fmt.Errorf("invalid message content, content: %s", content)
//	}
//
//	if !carbon.Parse(strTime).IsValid() {
//		return nil, fmt.Errorf("invalid message time, time: %s", strTime)
//	}
//
//	sender := content[:index]
//	msgContent := strings.TrimSpace(content[index+1:])
//
//	return &WechatMessage{
//		Created: strTime,
//		Scene:   scene,
//		Sender:  sender,
//		Content: msgContent,
//	}, nil
//}
