package parse

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type MockReaderWriter struct {
	data []byte
}

func (ns *MockReaderWriter) Read(p []byte) (int, error) {
	n := copy(p, ns.data)
	return n, nil
}

func (ns *MockReaderWriter) Write(p []byte) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}
	ns.data = p
	return len(ns.data), nil
}

func TestMultipleSpaces(t *testing.T) {
	rw := &MockReaderWriter{data: []byte("GET /?name=bob%20jones%20jimbob HTTP/1.1")}
	ParseRequest(rw)
	require.Equal(t, "HTTP/1.1 200 OK\r\nContent-Type: text/html; charset=utf-8\r\n"+
		"Content-Length: 31\r\n\r\n<h1>Hello bob jones jimbob</h1>", string(rw.data))
}

func TestNoSpace(t *testing.T) {
	rw := &MockReaderWriter{data: []byte("GET /?name=bob HTTP/1.1")}
	ParseRequest(rw)
	require.Equal(t, "HTTP/1.1 200 OK\r\nContent-Type: text/html; charset=utf-8\r\n"+
		"Content-Length: 18\r\n\r\n<h1>Hello bob</h1>", string(rw.data))
}

func TestMissingData(t *testing.T) {
	rw := &MockReaderWriter{data: []byte("GET / HTTP/1.1")}
	ParseRequest(rw)
	require.Equal(t, "HTTP/1.1 400 Bad Request\r\nContent-Type: text/plain; charset=utf-8\r\n"+
		"Content-Length: 29\r\n\r\nserver received a bad request", string(rw.data))
}

func TestBadData(t *testing.T) {
	rw := &MockReaderWriter{data: []byte("GET /?name=bob%jones HTTP/1.1")}
	ParseRequest(rw)
	require.Equal(t, "HTTP/1.1 400 Bad Request\r\nContent-Type: text/plain; charset=utf-8\r\n"+
		"Content-Length: 29\r\n\r\nserver received a bad request", string(rw.data))
}

func TestMethodNotAllowed(t *testing.T) {
	rw := &MockReaderWriter{data: []byte("POST /?name=bob%jones HTTP/1.1")}
	ParseRequest(rw)
	require.Equal(t, "HTTP/1.1 405 Method Not Allowed\r\n" +
		"Content-Type: text/plain; charset=utf-8\r\nContent-Length: 18\r\n\r\nmethod not allowed", string(rw.data))
}
