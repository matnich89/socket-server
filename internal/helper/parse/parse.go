package parse

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

var (
	ErrBadRequest     = errors.New("server received a bad request")
	ErrInternalServer = errors.New("server encountered an error")
)

func ParseRequest(rw io.ReadWriter) {
	bytes := make([]byte, 1024) // I've made the assumption 1024 is enough
	n, err := rw.Read(bytes)

	if err != nil {
		writeBadRequestResponse(rw)
	}

	request := string(bytes[:n])

	requestMethod := strings.Split(request, " ")[0]

	if requestMethod != http.MethodGet {
		writeMethodNotAllowedResponse(rw)
		return
	}

	var name string

	_, err = fmt.Sscanf(request, "GET /?name=%s HTTP/1.1", &name)
	if err != nil {
		writeBadRequestResponse(rw)
		return
	}

	decodedName, err := url.QueryUnescape(name)

	if err != nil {
		writeBadRequestResponse(rw)
		return
	}

	trimmedName := strings.TrimSpace(decodedName)

	responseStr := "<h1>Hello " + trimmedName + "</h1>"

	writeOkResponse(rw, responseStr)
	if err != nil {
		writeInternalServerResponse(rw)
	}
}

func writeOkResponse(rw io.ReadWriter, response string) {
	_, err := io.WriteString(rw, fmt.Sprintf("HTTP/1.1 200 OK\r\n"+
		"Content-Type: text/html; charset=utf-8\r\n"+
		"Content-Length: %d\r\n"+
		"\r\n"+
		"%s", len(response), response))

	if err != nil {
		writeInternalServerResponse(rw)
	}
}

func writeBadRequestResponse(rw io.ReadWriter) {
	_, _ = io.WriteString(rw, fmt.Sprintf("HTTP/1.1 400 Bad Request\r\n"+
		"Content-Type: text/plain; charset=utf-8\r\n"+
		"Content-Length: %d\r\n"+
		"\r\n"+
		"%s", len(ErrBadRequest.Error()), ErrBadRequest.Error()))
}

func writeInternalServerResponse(rw io.ReadWriter) {
	_, _ = io.WriteString(rw, fmt.Sprintf("HTTP/1.1 500 Internal Server Error\r\n"+
		"Content-Type: text/plain; charset=utf-8\r\n"+
		"Content-Length: %d\r\n"+
		"\r\n"+
		"%s", len(ErrInternalServer.Error()), ErrInternalServer.Error()))
}

func writeMethodNotAllowedResponse(rw io.ReadWriter) {
	message := "method not allowed"
	_, _ = io.WriteString(rw, fmt.Sprintf("HTTP/1.1 405 Method Not Allowed\r\n"+
		"Content-Type: text/plain; charset=utf-8\r\n"+
		"Content-Length: %d\r\n"+
		"\r\n"+
		"%s", len(message), message))
}
