package x

import (
	"encoding/binary"
	// "fmt"
	"reflect"
)

func (b *Backend) marshall(data interface{}) {
	b.write(data)
}

func (b *Backend) unmarshall(data interface{}) {
	b.unmarshallValue(reflect.ValueOf(data).Elem())
}

func (b *Backend) unmarshallValue(value reflect.Value) {
	switch value.Kind() {
	case reflect.Struct:
		for i := 0; i < value.NumField(); i++ {
			b.unmarshallField(value, i)
		}
	case reflect.Array:
		for i := 0; i < value.Len(); i++ {
			b.unmarshallValue(value.Index(i))
		}
	case reflect.Slice:
		panic("cannot unmashall value without length information")
	default:
		b.read(value.Addr().Interface())
	}
}

func (b *Backend) unmarshallField(value reflect.Value, field int) {
	fieldValue := value.Field(field)
	fieldKind := fieldValue.Kind()
	sfield := value.Type().Field(field)

	if fieldKind != reflect.Slice && fieldKind != reflect.String {
		b.unmarshallValue(fieldValue)
		return
	}

	fieldType := fieldValue.Type()
	fieldTag := sfield.Tag

	lengthField := fieldTag.Get("lengthField")
	if lengthField == "" {
		panic("no length field for " + sfield.Name)
	}

	lengthValue := value.FieldByName(lengthField)
	length := lengthValue.Uint()

	if fieldKind == reflect.Slice {
		slc := fieldValue
		tmp := reflect.New(fieldType.Elem()).Elem()

		for i := 0; i < int(length); i++ {
			b.unmarshallValue(tmp)
			slc = reflect.Append(slc, tmp)
		}

		fieldValue.Set(slc)
	} else {
		var buf []byte
		var tmp byte

		for i := 0; i < int(length); i++ {
			b.read(&tmp)
			buf = append(buf, tmp)
		}
		b.readPadding()
		fieldValue.SetString(string(buf))
	}
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
	b.readUnused((4 - b.bytesRead%4) % 4)
}
