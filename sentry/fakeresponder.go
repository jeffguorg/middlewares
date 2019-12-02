package sentry

import (
	"bytes"
	"net/http"
)

type FakeResponseWriter struct {
	buffer bytes.Buffer
	status int

	w http.ResponseWriter
}

func (writer *FakeResponseWriter) Write(data []byte) (int, error) {
	return writer.buffer.Write(data)
}

func (writer *FakeResponseWriter) Header() http.Header {
	return writer.w.Header()
}

func (writer *FakeResponseWriter) WriteHeader(status int) {
	writer.status = status
}

func (writer *FakeResponseWriter) Flush() {
	writer.w.WriteHeader(writer.status)
	_, _ = writer.w.Write(writer.buffer.Bytes())

}
