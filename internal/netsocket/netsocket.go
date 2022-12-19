package netsocket

import (
	"log"
	"net"
	"os"
	"syscall"

	"syscallserver/internal/helper/parse"
)

type NetSocket struct {
	fd int
}

func New(ipStr string, port int) (*NetSocket, error) {
	ip := net.ParseIP(ipStr)

	syscall.ForkLock.Lock()
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, syscall.IPPROTO_TCP)
	if err != nil {
		return nil, os.NewSyscallError("socket", err)
	}
	syscall.ForkLock.Unlock()

	if err = syscall.SetsockoptInt(fd, syscall.SOL_SOCKET, syscall.SO_REUSEADDR, syscall.IPPROTO_TCP); err != nil {
		return nil, os.NewSyscallError("setsockopt", err)
	}

	sa := &syscall.SockaddrInet4{
		Port: port,
	}
	copy(sa.Addr[:], ip)

	if err = syscall.Bind(fd, sa); err != nil {
		return nil, err
	}

	if err = syscall.Listen(fd, syscall.SOMAXCONN); err != nil {
		return nil, err
	}

	return &NetSocket{fd: fd}, nil
}

func (ns *NetSocket) Accept() (*NetSocket, error) {
	nfd, _, err := syscall.Accept(ns.fd)
	if err != nil {
		return nil, err
	}
	return &NetSocket{fd: nfd}, nil
}

func (ns *NetSocket) Close() error {
	return syscall.Close(ns.fd)
}

func (ns *NetSocket) Listen() {
	for {
		conn, err := ns.Accept()
		if err != nil {
			log.Println(err)
			err := ns.Close()
			if err != nil {
				return
			}
		}
		parse.ParseRequest(conn)
		_ = conn.Close()
	}
}

func (ns *NetSocket) Read(p []byte) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}
	n, err := syscall.Read(ns.fd, p)
	if err != nil {
		n = 0
	}
	return n, err
}

func (ns *NetSocket) Write(p []byte) (int, error) {
	n, err := syscall.Write(ns.fd, p)
	if err != nil {
		n = 0
	}
	return n, err
}
