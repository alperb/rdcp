package common

import (
	"fmt"
	"time"
)

type Response struct {
	Method    METHOD
	Headers   map[string]string
	Body      RequestBody
	Timestamp int64
}

func NewResponseFromString(s string) (*Response, error) {
	r := NewResponse(METHOD(s))
	return r, nil
}

func NewResponse(m METHOD) *Response {
	return &Response{
		Method:    m,
		Headers:   make(map[string]string),
		Body:      RequestBody{},
		Timestamp: int64(time.Now().Unix()),
	}
}

func (r *Response) AddHeader(key string, value string) {
	r.Headers[key] = value
}

func (r *Response) Write(b []byte) error {
	_, err := r.Body.Write(b)
	return err
}

func (r *Response) ReadBody() []byte {
	b := make([]byte, 1024)
	r.Body.Read(b)
	return b
}

func (r *Response) GetSerializedHeaders() string {
	var s string
	for k, v := range r.Headers {
		s += fmt.Sprintf("%s:%s\r\n", k, v)
	}
	if len(s) > 2 {
		return s[0 : len(s)-2]
	} else {
		return s
	}
}

func (r *Response) GetSerializedBody() string {
	return string(r.ReadBody())
}

func (r *Response) GetSerialized() []byte {
	return []byte(fmt.Sprintf("%s\r\n%s\r\n\r\n%s\r\n", r.Method, r.GetSerializedHeaders(), r.GetSerializedBody()))
}

func (r *Response) ToWriteable() []byte {
	return r.GetSerialized()
}
