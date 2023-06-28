package common

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type Request struct {
	Method    METHOD
	Headers   map[string]string
	Body      RequestBody
	Timestamp int64
}

func CreateDefaultHeaders(ts time.Time, v string, response bool) map[string]string {
	if response {
		return map[string]string{
			"Version":  v,
			"Received": fmt.Sprintf("%d", ts.Unix()),
			"Replied":  fmt.Sprintf("%d", time.Now().Unix()),
		}
	} else {
		return map[string]string{
			"Time":    fmt.Sprintf("%d", ts.Unix()),
			"Version": v,
		}
	}
}

func NewRequest(method METHOD) *Request {
	return &Request{
		Method:    method,
		Headers:   make(map[string]string),
		Body:      RequestBody{},
		Timestamp: int64(time.Now().Unix()),
	}
}

func (r *Request) AddHeader(key string, value string) {
	r.Headers[key] = value
}

func (r *Request) WriteBody(c []byte) {
	r.Body.Write(c)
}

func (r *Request) WriteBodyString(s string) {
	r.Body.Write([]byte(s))
}

func (r *Request) ReadBody() []byte {
	b := make([]byte, 1024)
	r.Body.Read(b)
	return b
}

func (r *Request) GetSerializedHeaders() string {
	var s string
	for k, v := range r.Headers {
		s += fmt.Sprintf("%s:%s\r\n", k, v)
	}
	return s
}

func (r *Request) GetActionName() string {
	if r.Method != ORDER {
		return ""
	}

	b := make([]byte, 1024)
	r.Body.Read(b)
	s := string(b)

	splitted := strings.Split(s, "(")
	if len(splitted) < 2 {
		return ""
	}

	return splitted[0]
}

func (r *Request) GetActionParameters() map[string]string {
	if r.Method != ORDER {
		return nil
	}

	b := make([]byte, 1024)
	r.Body.Read(b)
	s := string(b)

	splitted := strings.Split(s, "(")
	if len(splitted) < 2 {
		return nil
	}
	if len(splitted[1]) == 1 {
		return make(map[string]string)
	}

	paramsTemp := strings.Split(splitted[1], ")")
	params := paramsTemp[0]
	paramsSplitted := strings.Split(params, ",")

	m := make(map[string]string)
	for _, param := range paramsSplitted {
		split := strings.Split(param, "=")
		if len(split) < 2 {
			continue
		}
		m[strings.Trim(split[0], " ")] = strings.Trim(split[1], " ")
	}
	return m
}

func (r *Request) GetSerializedBody() string {
	return string(r.ReadBody())
}

func (r *Request) GetSerialized() []byte {
	return []byte(fmt.Sprint("%s\r\n%s\r\n\r\n%s", r.Method, r.GetSerializedHeaders(), r.GetSerializedBody()))
}

func (r *Request) Send(host string) *Response {
	res := &Response{}
	return res
}

func ParseContentIntoRequest(content []byte) *Request {
	s := string(content)
	split := strings.Split(s, "\\r\\n\\r\\n")

	var method METHOD
	var headers map[string]string
	var body RequestBody
	var timestamp int64

	if len(split) < 2 {
		return &Request{}
	}

	uppersplit := strings.Split(split[0], "\\r\\n")

	if len(split) > 0 {
		method = ParseMethod(uppersplit[0])
	}

	if len(split) > 1 {
		headers, timestamp = ParseHeaders(uppersplit)
	}

	body.Write([]byte(split[1]))

	return &Request{
		Method:    method,
		Headers:   headers,
		Body:      body,
		Timestamp: timestamp,
	}
}

func ParseHeaders(content []string) (map[string]string, int64) {
	headers := make(map[string]string)
	var timestamp int64

	for idx, line := range content {
		if idx == 0 {
			continue
		}

		if line == "" {
			break
		}

		if strings.HasPrefix(line, "Timestamp") {
			ts, e := strconv.ParseInt(strings.Split(line, ":")[1], 10, 64)
			if e != nil {
				log.Fatal(e)
				os.Exit(1)
			}
			timestamp = ts
		} else {
			split := strings.Split(line, ":")
			if len(split) > 1 {
				headers[split[0]] = split[1]
			}
		}
	}

	return headers, timestamp
}

func ParseMethod(content string) METHOD {
	switch content {
	case "ORD":
		return ORDER
	case "R/ORD":
		return ORDER_RESPONSE
	case "PRX":
		return PROXY
	case "R/PRX":
		return PROXY_RESPONSE
	case "OPT":
		return OPTIONS
	case "R/OPT":
		return OPTIONS_RESPONSE
	case "HBT":
		return HEARTBEAT
	case "R/HBT":
		return HEARTBEAT_RESPONSE
	}
	return ""
}
