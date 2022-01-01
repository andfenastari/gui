package wayland

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"reflect"
)

type Message struct {
	ObjectId ObjectId
	Size     uint16
	Opcode   Opcode
	Payload  []byte
}

func NewMessage(objectId ObjectId, opcode Opcode, data interface{}) (msg Message) {
	var b bytes.Buffer
	err := marshall(&b, data)
	if err != nil {
		panic(err)
	}
	return NewMessageBytes(objectId, opcode, b.Bytes())
}

func NewMessageBytes(objectId ObjectId, opcode Opcode, payload []byte) (msg Message) {
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

func unmarshall(r io.Reader, data interface{}) (err error) {
	dataValue := reflect.ValueOf(data)
	if dataValue.Kind() != reflect.Ptr || dataValue.Elem().Kind() != reflect.Struct {
		panic("'data' must be a pointer to a struct value")
	}

	elemValue := dataValue.Elem()
	numField := elemValue.NumField()

	for i := 0; i < numField; i++ {
		field := elemValue.Field(i)

		switch field.Kind() {
		case reflect.Uint32:
			var n uint32
			n, err = unmarshallUint32(r)
			field.SetUint(uint64(n))
		case reflect.Int32:
			var n int32
			n, err = unmarshallInt32(r)
			field.SetInt(int64(n))
		case reflect.Float32:
			var n float32
			n, err = unmarshallFloat32(r)
			field.SetFloat(float64(n))
		case reflect.String:
			var s string
			s, err = unmarshallString(r)
			field.SetString(s)
		case reflect.Array:
			var a []byte
			if field.Elem().Kind() != reflect.Uint8 {
				unsuportedType(field)
			}
			a, err = unmarshallArray(r)
			field.SetBytes(a)
		default:
			unsuportedType(field)
		}
		if err != nil {
			return fmt.Errorf("unmarshalling field '%s': %w",
				elemValue.Type().Field(i).Name,
				err,
			)
		}
	}
	return nil
}

func unsuportedType(field reflect.Value) {
	panic("Unsuported field type '" + field.Type().Name() + "'")
}

func unmarshallUint32(r io.Reader) (v uint32, err error) {
	var buf [4]byte
	_, err = r.Read(buf[:])
	if err != nil {
		return v, fmt.Errorf("unmarshalling uint32 value: %w", err)
	}
	return binary.LittleEndian.Uint32(buf[:]), nil
}

func unmarshallInt32(r io.Reader) (v int32, err error) {
	var buf [4]byte
	_, err = r.Read(buf[:])
	if err != nil {
		return v, fmt.Errorf("unmarshalling int32 value: %w", err)
	}
	return int32(binary.LittleEndian.Uint32(buf[:])), nil
}

func unmarshallFloat32(r io.Reader) (v float32, err error) {
	var buf [4]byte
	_, err = r.Read(buf[:])
	if err != nil {
		return v, fmt.Errorf("unmarshalling float32 value: %w", err)
	}
	b := binary.LittleEndian.Uint32(buf[:])
	return math.Float32frombits(b), nil
}

func unmarshallString(r io.Reader) (s string, err error) {
	var buf [4]byte
	_, err = r.Read(buf[:])
	if err != nil {
		return s, fmt.Errorf("unmarshalling string value: %w", err)
	}

	l := int(binary.LittleEndian.Uint32(buf[:]))
	p := (4 - l%4) % 4
	b := make([]byte, l+p)
	_, err = r.Read(b)
	if err != nil {
		return s, fmt.Errorf("unmarshalling string value: %w", err)
	}

	return string(b[0 : l-1]), nil
}

func unmarshallArray(r io.Reader) (b []byte, err error) {
	var buf [4]byte
	_, err = r.Read(buf[:])
	if err != nil {
		return b, fmt.Errorf("unmarshalling array value: %w", err)
	}

	l := int(binary.LittleEndian.Uint32(buf[:]))
	p := (4 - l%4) % 4
	b = make([]byte, l+p)
	_, err = r.Read(b)
	if err != nil {
		return b, fmt.Errorf("unmarshalling array value: %w", err)
	}

	return b[0:l], nil
}

func marshall(w io.Writer, data interface{}) (err error) {
	elemValue := reflect.ValueOf(data)
	if elemValue.Kind() != reflect.Struct {
		panic("'data' must be a struct value")
	}

	numField := elemValue.NumField()

	for i := 0; i < numField; i++ {
		field := elemValue.Field(i)

		switch field.Kind() {
		case reflect.Uint32:
			err = marshallUint32(w, uint32(field.Uint()))
		case reflect.Int32:
			err = marshallInt32(w, int32(field.Int()))
		case reflect.Float32:
			err = marshallFloat32(w, float32(field.Float()))
		case reflect.String:
			err = marshallString(w, field.String())
		case reflect.Array:
			if field.Elem().Kind() != reflect.Uint8 {
				unsuportedType(field)
			}
			err = marshallArray(w, field.Bytes())
		default:
			unsuportedType(field)
		}
		if err != nil {
			return fmt.Errorf("marshalling field '%s': %w",
				elemValue.Type().Field(i).Name,
				err,
			)
		}
	}
	return nil
}

func marshallUint32(w io.Writer, v uint32) (err error) {
	var buf [4]byte
	binary.LittleEndian.PutUint32(buf[:], v)
	_, err = w.Write(buf[:])
	if err != nil {
		return fmt.Errorf("marshaling uint32 value: %w", err)
	}
	return nil
}

func marshallInt32(w io.Writer, v int32) (err error) {
	var buf [4]byte
	binary.LittleEndian.PutUint32(buf[:], uint32(v))
	_, err = w.Write(buf[:])
	if err != nil {
		return fmt.Errorf("marshaling int32 value: %w", err)
	}
	return nil
}

func marshallFloat32(w io.Writer, v float32) (err error) {
	n := math.Float32bits(v)
	var buf [4]byte
	binary.LittleEndian.PutUint32(buf[:], n)
	_, err = w.Write(buf[:])
	if err != nil {
		return fmt.Errorf("marshalling float32 value: %w", err)
	}
	return nil
}

func marshallString(w io.Writer, s string) (err error) {
	l := len(s) + 1
	p := (4 - l%4) % 4

	var buf [4]byte
	binary.LittleEndian.PutUint32(buf[:], uint32(l))
	_, err = w.Write(buf[:])
	if err != nil {
		return fmt.Errorf("marshalling string value: %w", err)
	}

	b := make([]byte, l+p+1)
	copy(b, []byte(s))
	_, err = w.Write(b)
	if err != nil {
		return fmt.Errorf("marshalling string value: %w", err)
	}
	return nil
}

func marshallArray(w io.Writer, a []byte) (err error) {
	l := len(a)
	p := (4 - l%4) % 4

	var buf [4]byte
	binary.LittleEndian.PutUint32(buf[:], uint32(l))
	_, err = w.Write(buf[:])
	if err != nil {
		return fmt.Errorf("marshalling array value: %w", err)
	}

	b := make([]byte, l+p)
	copy(b, a)
	_, err = w.Write(b)
	if err != nil {
		return fmt.Errorf("marshalling array value: %w", err)
	}
	return err
}
