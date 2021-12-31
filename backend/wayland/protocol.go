package wayland

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"sync/atomic"
)

type ObjectId uint32
type Opcode uint16

var prevObjectId uint32 = 1

const DisplayId ObjectId = 1

const (
	OpGetRegistry Opcode = 1
)

func NewObjectId() ObjectId {
	return ObjectId(atomic.AddUint32(&prevObjectId, 1))
}

type Message struct {
	ObjectId ObjectId
	Size     uint16
	Opcode   Opcode
	Payload  []byte
}

type EventRegistryGlobal struct {
	Name      uint32
	Interface string
	Version   uint32
}

func NewMessage(objectId ObjectId, opcode Opcode, payload []byte) (msg Message) {
	msg.Size = uint16(8 + len(payload))
	msg.Opcode = opcode
	msg.ObjectId = objectId
	msg.Payload = payload
	return
}

func ReadMessage(r io.Reader) (msg Message, err error) {
	var header [8]byte
	_, err = r.Read(header[:])
	if err != nil {
		return msg, fmt.Errorf("reading message header: %w", err)
	}

	sizeAndOpcode := binary.LittleEndian.Uint32(header[4:8])
	msg.ObjectId = ObjectId(binary.LittleEndian.Uint32(header[0:4]))
	msg.Size = uint16(sizeAndOpcode >> 16)
	msg.Opcode = Opcode(sizeAndOpcode)

	msg.Payload = make([]byte, msg.Size-8)
	_, err = r.Read(msg.Payload)
	if err != nil {
		return msg, fmt.Errorf("reading message payload: %w", err)
	}

	return
}

func (msg *Message) Write(w io.Writer) (err error) {
	var header [8]byte
	sizeAndOpcode := (uint32(msg.Size) << 16) | uint32(msg.Opcode)
	binary.LittleEndian.PutUint32(header[0:4], uint32(msg.ObjectId))
	binary.LittleEndian.PutUint32(header[4:8], sizeAndOpcode)

	_, err = w.Write(header[:])
	if err != nil {
		return fmt.Errorf("writing message header: %w", err)
	}

	_, err = w.Write(msg.Payload)
	if err != nil {
		return fmt.Errorf("writing message payload: %w", err)
	}

	return nil
}

func (msg *Message) Unmarshall(data interface{}) {
	buf := bytes.NewBuffer(msg.Payload)
	err := unmarshall(buf, data)
	if err != nil {
		panic(err)
	}
}

func NewGetRegistry(registryId ObjectId) (msg Message) {
	payload := make([]byte, 4)
	binary.LittleEndian.PutUint32(payload, uint32(registryId))
	return NewMessage(
		DisplayId,
		OpGetRegistry,
		payload,
	)
}
