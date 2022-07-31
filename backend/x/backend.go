package x

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"log"
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
	nextId       Card32
}

type Window struct{}

func (b *Backend) Init() (err error) {
	b.conn, err = connect()
	if err != nil {
		return fmt.Errorf("initializing backend connection: %w", err)
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
		return fmt.Errorf("sending init request: %w", b.err)
	}

	var success Card8
	b.read(&success)

	if b.err != nil {
		return fmt.Errorf("reading init response: %w", b.err)
	}

	if success != 1 {
		return fmt.Errorf("init response %d: %w", success, ErrNotImplemented)
	}

	//b.readInitResponse(&b.initResponse)
	b.unmarshall(&b.initResponse)
	if b.err != nil {
		return fmt.Errorf("reading init response: %w", b.err)
	}

	pprint(b.initResponse)

	return nil
}

func (b *Backend) Close() {
	b.conn.Close()
}

func (b *Backend) OpenWindow(title string, width, height int) (w *Window, err error) {
	w = new(Window)

	return w, nil
}

func (b *Backend) allocId() (n Card32) {
	n = b.nextId | b.initResponse.ResourceIdBase
	b.nextId++
	return
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

func pprint(v interface{}) {
	j, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(j))
}
