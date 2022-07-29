package wayland

import (
	"fmt"
	"log"
	"net"
	"os"
	"path"
	"strconv"
	"sync/atomic"
)

type Backend struct {
	Conn         net.Conn
	File         *os.File
	PrevObjectId uint32

	RegistryId   ObjectId
	CompositorId ObjectId
	SurfaceId    ObjectId
}

func NewBackend() (b Backend, err error) {
	b.Conn, b.File, err = connect()
	if err != nil {
		return b, fmt.Errorf("initializing wayland socket connection: %w", err)
	}

	b.PrevObjectId = 1
	b.RegistryId = b.NewObjectId()
	syncCallbackId := b.NewObjectId()

	var msg Message
	msg = NewMessage(DisplayId, OpDisplayGetRegistry, DisplayGetRegistry{
		Registry: b.RegistryId,
	})

	err = msg.Write(b.Conn)
	if err != nil {
		return b, fmt.Errorf("writing get_registry message: %w", err)
	}

	msg = NewMessage(DisplayId, OpDisplaySync, DisplaySync{
		Callback: syncCallbackId,
	})
	err = msg.Write(b.Conn)
	if err != nil {
		return b, fmt.Errorf("writing sync message: %w", err)
	}

	done := false
	for !done {
		msg, err = ReadMessage(b.Conn)
		if err != nil {
			return b, fmt.Errorf("reading registry global message: %w", err)
		}
		log.Printf("received message: %v", msg)

		switch {
		case msg.ObjectId == syncCallbackId:
			log.Printf("received callback %d", msg.ObjectId)
			done = true

		case msg.ObjectId == b.RegistryId && msg.Opcode == OpRegistryGlobal:
			var ev RegistryGlobal
			msg.Unmarshall(&ev)
			log.Printf("global: %v", ev)

			if ev.Interface == "wl_compositor" {
				b.CompositorId = b.NewObjectId()
				msg = NewMessage(b.RegistryId, OpRegistryBind, RegistryBind{
					Name:      ev.Name,
					Interface: ev.Interface,
					Version:   ev.Version,
					Id:        b.CompositorId,
				})
				log.Printf("bind message: %v", msg)
				err = msg.Write(b.Conn)
				if err != nil {
					return b, fmt.Errorf("binding to compositor: %w", err)
				}
				log.Printf("bound to compositor %d", b.CompositorId)
			}
		case msg.ObjectId == DisplayId && msg.Opcode == OpDisplayError:
			var ev DisplayError
			msg.Unmarshall(&ev)
			log.Fatalf("error: %v", ev)
		}

	}

	b.SurfaceId = b.NewObjectId()
	msg = NewMessage(b.CompositorId, OpCompositorCreateSurface, CompositorCreateSurface{
		Id: b.SurfaceId,
	})
	err = msg.Write(b.Conn)
	if err != nil {
		return b, fmt.Errorf("creating suface: %w", err)
	}
	log.Printf("created surface %d", b.SurfaceId)

	return
}

func (b *Backend) Close() {
	b.Conn.Close()
	if b.File != nil {
		b.File.Close()
	}
}

func (b *Backend) NewObjectId() ObjectId {
	return ObjectId(atomic.AddUint32(&b.PrevObjectId, 1))
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
