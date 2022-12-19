package netsocket

import (
	"io"
	"log"
	"net/http"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/require"
)

/*
The go routines and
multiple requests might seem
crazy, but I feel it gives
confidence the server is able to handle
multiple requests

128 is the max number of connections that syscall.SOMAXCONN is set to
on my machine
*/
func TestServerOK(t *testing.T) {
	startServer(t)
	client := http.Client{}

	var successfulRequests int32
	wg := sync.WaitGroup{}
	for i := 0; i < 128; i++ {
		wg.Add(1)
		go func(successfulRequests *int32, wg *sync.WaitGroup) {
			req, err := http.NewRequest(http.MethodGet, "http://localhost:8889/?name=bob%20jones", nil)
			require.NoError(t, err)
			resp, err := client.Do(req)
			require.NoError(t, err)
			b, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			require.Equal(t, "200 OK", resp.Status)
			require.Equal(t, "text/html; charset=utf-8", resp.Header.Get("Content-Type"))
			require.Equal(t, "<h1>Hello bob jones</h1>", string(b))
			log.Println(string(b))
			atomic.AddInt32(successfulRequests, 1)
			wg.Done()
		}(&successfulRequests, &wg)
	}
	wg.Wait()
	require.Equal(t, int32(128), successfulRequests)
}

func startServer(t *testing.T) {
	t.Helper()
	server, err := New("127.0.0.1", 8889)
	require.NoError(t, err)
	require.NoError(t, err)
	go server.Listen()
}
