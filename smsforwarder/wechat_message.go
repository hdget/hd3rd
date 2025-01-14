package smsforwarder

import (
	"fmt"
	"github.com/golang-module/carbon/v2"
	"github.com/pkg/errors"
	"regexp"
	"strings"
)

type WechatMessageKind int

const (
	WechatMessageKindUnknown WechatMessageKind = iota
	WechatMessageKindOne2One                   // 单聊消息
	WechatMessageKindGroup                     // 群聊消息
)

type WechatMessage struct {
	Kind    WechatMessageKind
	Created string
	Group   string
	Sender  string
	Content string
}

var (
	regexRemoveUnused = regexp.MustCompile(`^\[[^\]]*\]`)
)

func NewWechatMessage(data []byte, secret string) (*WechatMessage, error) {
	rawMessage, err := NewRawMessage(data, secret)
	if err != nil {
		return nil, err
	}

	return parseContent(rawMessage.Content)
}

func parseContent(s string) (*WechatMessage, error) {
	if s == "" {
		return nil, errors.New("empty wechat message")
	}

	tokens := strings.Split(s, "\n")
	tokenLength := len(tokens)
	if tokenLength < 6 {
		return nil, fmt.Errorf("invalid wechat message, rawContent: %s", s)
	}

	strTime, scene, content := tokens[tokenLength-2], tokens[tokenLength-4], regexRemoveUnused.ReplaceAllString(strings.Join(tokens[1:tokenLength-4], "\n"), "")

	index := strings.Index(content, ":")
	if index == -1 {
		return nil, fmt.Errorf("invalid wechat message content, content: %s", content)
	}

	if !carbon.Parse(strTime).IsValid() {
		return nil, fmt.Errorf("invalid wechat message time, time: %s", strTime)
	}

	// 如果消息的联系人和发送者是一样，则表示是单聊消息，否则认为是群聊消息
	kind := WechatMessageKindOne2One
	group := ""
	sender := content[:index]
	if scene != sender {
		kind = WechatMessageKindGroup
		group = scene
	}

	return &WechatMessage{
		Kind:    kind,
		Created: strTime,
		Group:   group,
		Sender:  sender,
		Content: strings.TrimSpace(content[index+1:]),
	}, nil
}
