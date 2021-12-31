package wayland

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"reflect"
)

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
			field.Set(reflect.ValueOf(a))
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

	return string(b[0:l]), nil
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
