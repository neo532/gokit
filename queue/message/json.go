package message

import "encoding/json"

/*
 * @abstract json message
 * @mail neo532@126.com
 * @date 2023-08-13
 */

type Json[T any] struct {
	MsgID string `json:"msgId"`
	Tag   string `json:"tag"`
	Data  T      `json:"data"`
}

func NewJson[T any]() *Json[T] {
	return &Json[T]{}
}

func (m *Json[T]) Marshal(msgID string, tag string, data T) (b []byte, err error) {
	m.MsgID = msgID
	m.Data = data
	m.Tag = tag
	b, err = json.Marshal(m)
	return
}

func (m *Json[T]) Data() T {
	return m.Data
}

func (m *Json[T]) Unmarshal(data []byte) error {
	return json.Unmarshal(data, m)
}
