package common

import (
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

type Request struct {
	Method    METHOD
	Headers   map[string]string
	Body      io.Reader
	Timestamp int64
}

func ParseContentIntoRequest(content []byte) *Request {
	s := string(content)
	split := strings.Split(s, "\r\n\r\n")

	var method METHOD
	var headers map[string]string
	var body io.Reader
	var timestamp int64

	if len(split) > 0 {
		method = ParseMethod(split[0])
	}

	if len(split) > 1 {
		headers, timestamp = ParseHeaders(split)
	}

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

type Response struct {
	Method    METHOD
	Headers   map[string]string
	Body      io.Writer
	Timestamp int64
}

func (r *Response) Write(b []byte) error {
	_, err := r.Body.Write(b)
	return err
}
