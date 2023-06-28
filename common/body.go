package common

type RequestBody struct {
	Content []byte
}

func (rb *RequestBody) Read(p []byte) (n int, err error) {
	copy(p, rb.Content)
	return len(rb.Content), nil
}

func (rb *RequestBody) Write(p []byte) (n int, err error) {
	rb.Content = p
	return len(p), nil
}
