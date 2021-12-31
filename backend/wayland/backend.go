package wayland

import (
	"fmt"
	"log"
	"net"
	"os"
	"path"
	"strconv"
)

type Backend struct {
	Conn         net.Conn
	File         *os.File
	PrevObjectId uint32
}

func NewBackend() (b Backend, err error) {
	b.Conn, b.File, err = connect()
	b.PrevObjectId = 1

	var msg Message
	msg = NewGetRegistry(NewObjectId())
	err = msg.Write(b.Conn)
	if err != nil {
		return b, fmt.Errorf("writing get_registry message: %w", err)
	}
	log.Print(msg)
	for {
		msg, err = ReadMessage(b.Conn)
		if err != nil {
			return b, fmt.Errorf("reading registry global message: %w", err)
		}

		var ev EventRegistryGlobal
		msg.Unmarshall(&ev)
		log.Printf("%[1]T %[1]v", ev)	
	}

	return
}

func (b *Backend) Close() {
	b.Conn.Close()
	b.File.Close()
}

func connect() (conn net.Conn, f *os.File, err error) {
	socketFd := os.Getenv("WAYLAND_SOCKET")
	if socketFd != "" {
		socketFdI, err := strconv.Atoi(socketFd)
		if err != nil {
			return conn, f, fmt.Errorf("parsing 'WAYLAND_SOCKET' value: %w", err)
		}

		f = os.NewFile(uintptr(socketFdI), "wayland-0")
		if f != nil {
			return conn, f, fmt.Errorf("interpreting 'WAYLAND_SOCKET as a file")
		}
		defer func() {
			if err != nil {
				f.Close()
			}
		}()

		conn, err = net.FileConn(f)
		if err != nil {
			return conn, f, fmt.Errorf("interpreting 'WAYLAND_SOCKET' as a socket: %w", err)
		}
		defer func() {
			if err != nil {
				conn.Close()
			}
		}()

		return conn, f, nil
	}

	wlDisplay := os.Getenv("WAYLAND_DISPLAY")
	if wlDisplay == "" {
		wlDisplay = "wayland-0"
	}
	if !path.IsAbs(wlDisplay) {
		xdgRtd := os.Getenv("XDG_RUNTIME_DIR")
		if xdgRtd == "" {
			return conn, f, fmt.Errorf("'WAYLAND_DISPLAY' is a relative path but 'XDG_RUNTIME_DIR' is not set")
		}

		wlDisplay = path.Join(xdgRtd, wlDisplay)
	}

	conn, err = net.Dial("unix", wlDisplay)
	if err != nil {
		return conn, f, fmt.Errorf("connecting to wayland display socket '%s': %w", wlDisplay, err)
	}

	return conn, nil, nil

}
