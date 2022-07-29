package x

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

var (
	ErrInit           = errors.New("initializing X connection")
	ErrNotImplemented = errors.New("Not implemented")
)

type Backend struct {
	conn         net.Conn
	byteOrder    binary.ByteOrder
	bytesWritten int
	bytesRead    int
	err          error

	initResponse InitResponse
}

func NewBackend() (b *Backend, err error) {
	b = new(Backend)
	b.conn, err = connect()
	if err != nil {
		return nil, fmt.Errorf("initializing backend connection: %w", err)
	}

	b.byteOrder = binary.BigEndian

	b.write(Card8('B')) // Big Endian
	b.writeUnused(1)
	b.write(Card16(11)) // Protocol major version
	b.write(Card16(0))  // Protocol minor version

	b.write(Card16(0)) // Auth protocol length
	b.write(Card16(0)) // Auth data length
	b.writeUnused(2)   // Padding

	if b.err != nil {
		return nil, fmt.Errorf("sending init request: %w", b.err)
	}

	var success Card8
	b.read(&success)

	if b.err != nil {
		return nil, fmt.Errorf("reading init response: %w", b.err)
	}

	if success != 1 {
		return nil, fmt.Errorf("init response %d: %w", success, ErrNotImplemented)
	}

	var ir InitResponse
	b.readInitResponse(&ir)
	if b.err != nil {
		return nil, fmt.Errorf("reading init response: %w", b.err)
	}
	
	fmt.Printf("Init Response = %#v\n", ir)

	return
}

func (b *Backend) Close() {
	b.conn.Close()
}

func (b *Backend) write(data interface{}) {
	if b.err != nil {
		return
	}

	b.err = binary.Write(b.conn, b.byteOrder, data)
	b.bytesWritten += binary.Size(data)
}

func (b *Backend) writeUnused(n int) {
	var buf [6]Card8
	b.write(buf[0:n])
}

func (b *Backend) writePadding() {
	b.writeUnused((4 - b.bytesWritten%4) % 4)
}

func (b *Backend) read(data interface{}) {
	if b.err != nil {
		return
	}

	b.err = binary.Read(b.conn, b.byteOrder, data)
	b.bytesRead += binary.Size(data)
}

func (b *Backend) readUnused(n int) {
	var buf [6]Card8
	b.read(buf[0:n])
}

func (b *Backend) readPadding() {
	b.readUnused((4-b.bytesRead%4)%4)
}

func connect() (conn net.Conn, err error) {
	host, display, _, err := parseDisplay(os.Getenv("DISPLAY"))
	if err != nil {
		err = fmt.Errorf("parsing `DISPLAY` environment variable: %w", err)
		return
	}

	if host == "" || host == "host/unix" {
		path := fmt.Sprintf("/tmp/.X11-unix/X%d", display)
		conn, err = net.Dial("unix", path)
	} else {
		port := 6000 + display
		path := fmt.Sprintf("host:%d", port)
		conn, err = net.Dial("tcp", path)
	}

	if err != nil {
		err = fmt.Errorf("connecting to X display: %w", err)
		return
	}

	return
}

func parseDisplay(spec string) (host string, display, screen int, err error) {
	host, spec, ok := strings.Cut(spec, ":")
	if !ok {
		err = fmt.Errorf("unspecified display: %w", ErrInit)
		return
	}

	displayS, screenS, _ := strings.Cut(spec, ".")

	if displayS == "" {
		err = fmt.Errorf("unspecified display: %w", ErrInit)
		return
	}

	display, err = strconv.Atoi(displayS)
	if err != nil {
		err = fmt.Errorf("parsing display number `%s`: %w", displayS, err)
		return
	}

	if screenS == "" {
		screen = 0
	} else {
		screen, err = strconv.Atoi(screenS)
		if err != nil {
			err = fmt.Errorf("parsing screen number `%s`: %w", screenS, err)
			return
		}
	}

	return
}
