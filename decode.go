package peinfo

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"unsafe"

	"golang.org/x/text/encoding/unicode"
)

var ErrEndofstream = errors.New("read end of stream. EOF")
var ErrLenOutofrange = errors.New("read len range of left length")

const (
	Lenint32 = 4
	Lenint64 = 8
	Lenlen   = 2
)

type ByteCode interface {
	Read(index int, stream []byte) (int, error)
	SetValue(value []byte)
}

func byteRead(index, llen int, stream []byte) ([]byte, int, error) {
	start := index
	end := index + llen

	if len(stream[start:]) < end-start {
		return nil, 0, ErrLenOutofrange
	}

	value := stream[start:end]

	return value, end, nil
}

type WinWord struct {
	Data int32
}

func (b *WinWord) Read(index int, stream []byte) (int, error) {
	v, next, err := byteRead(index, 2, stream)

	b.SetValue(v)

	return next, err
}

func (b *WinWord) SetValue(value []byte) {
	b.Data = int32(Bytes2Int(value))
}

// int32 from byte stream
var _ ByteCode = (*WinInt32)(nil)

type WinInt32 struct {
	Data int32
}

func (b *WinInt32) Read(index int, stream []byte) (int, error) {
	v, next, err := byteRead(index, Lenint32, stream)

	b.SetValue(v)

	return next, err
}

func (b *WinInt32) SetValue(value []byte) {
	b.Data = int32(Bytes2Int(value))
}

// int64 from byte stream
var _ ByteCode = (*WinInt64)(nil)

type WinInt64 struct {
	Data int64
}

func (b *WinInt64) Read(index int, stream []byte) (int, error) {
	v, next, err := byteRead(index, Lenint64, stream)

	b.SetValue(v)

	return next, err
}

func (b *WinInt64) SetValue(value []byte) {
	b.Data = int64(Bytes2Int(value))
}

// ansi from byte stream
var _ ByteCode = (*WinAnsi)(nil)

type WinAnsi struct {
	Data string
}

func (b *WinAnsi) Read(index int, stream []byte) (int, error) {
	lenbyte, next, err := byteRead(index, Lenlen, stream)
	if err != nil {
		return 0, err
	}

	lenint := Bytes2Int(lenbyte)

	strbyte, next, err := byteRead(next, lenint, stream)

	b.SetValue(strbyte)

	return next, err
}

func (b *WinAnsi) SetValue(value []byte) {
	b.Data = Bytes2String(value)
}

// unicode from byte stream
var _ ByteCode = (*WinUnicode)(nil)

type WinUnicode struct {
	Data string
}

func (b *WinUnicode) Read(index int, stream []byte) (int, error) {
	lenbyte, next, err := byteRead(index, Lenlen, stream)
	if err != nil {
		return 0, err
	}

	lenint := Bytes2Int(lenbyte)

	strbyte, next, err := byteRead(next, lenint, stream)

	b.SetValue(strbyte)

	return next, err
}

func (b *WinUnicode) SetValue(value []byte) {
	b.Data = UnicodeBytes2String(value)
}

// byte from byte stream
var _ ByteCode = (*WinByte)(nil)

type WinByte struct {
	Data []byte
}

func (b *WinByte) Read(index int, stream []byte) (int, error) {
	lenbyte, next, err := byteRead(index, Lenlen, stream)
	if err != nil {
		return 0, err
	}

	lenint := Bytes2Int(lenbyte)

	strbyte, next, err := byteRead(next, lenint, stream)

	b.SetValue(strbyte)

	return next, err
}

func (b *WinByte) SetValue(value []byte) {
	b.Data = value
}

func DecodeMessage(stream []byte, obj interface{}) error {
	var index int = 8
	var err error

	v := reflect.ValueOf(obj).Elem()

	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i)
		tag := field.Tag
		type1 := tag.Get("type")

		switch type1 {
		case "int32":
			winint32 := new(WinInt32)
			index, err = winint32.Read(index, stream)
			v.FieldByName(field.Name).Set(reflect.ValueOf(winint32.Data))
		case "int32_x":
			winint32 := new(WinInt32)
			index, err = winint32.Read(index, stream)
			value := fmt.Sprintf("0x%s", strconv.FormatInt(int64(uint32(winint32.Data)), 16))
			v.FieldByName(field.Name).Set(reflect.ValueOf(value))
		case "int64":
			winint64 := new(WinInt64)
			index, err = winint64.Read(index, stream)
			v.FieldByName(field.Name).Set(reflect.ValueOf(winint64.Data))
		case "unicode":
			winunicode := new(WinUnicode)
			index, err = winunicode.Read(index, stream)
			v.FieldByName(field.Name).Set(reflect.ValueOf(winunicode.Data))
		case "ansi":
			winansi := new(WinAnsi)
			index, err = winansi.Read(index, stream)
			v.FieldByName(field.Name).Set(reflect.ValueOf(winansi.Data))
		case "byte_int32":
			winbyte := new(WinByte)
			index, err = winbyte.Read(index, stream)
			v.FieldByName(field.Name).Set(reflect.ValueOf(int32(Bytes2Int(winbyte.Data))))
		case "byte_x":
			winbyte := new(WinByte)
			index, err = winbyte.Read(index, stream)
			value := fmt.Sprintf("0x%s", strconv.FormatInt(int64(Bytes2Int(winbyte.Data)), 16))
			v.FieldByName(field.Name).Set(reflect.ValueOf(value))

		}
	}

	return err
}

func Bytes2Int(b []byte) int {
	if len(b) == 3 {
		b = append([]byte{0}, b...)
	}
	bytesBuffer := bytes.NewBuffer(b)
	switch len(b) {
	case 1:
		var tmp uint8
		binary.Read(bytesBuffer, binary.LittleEndian, &tmp)
		return int(tmp)
	case 2:
		var tmp uint16
		binary.Read(bytesBuffer, binary.LittleEndian, &tmp)
		return int(tmp)
	case 4:
		var tmp uint32
		binary.Read(bytesBuffer, binary.LittleEndian, &tmp)
		return int(tmp)
	case 8:
		var tmp uint64
		binary.Read(bytesBuffer, binary.LittleEndian, &tmp)
		return int(tmp)
	default:
		return 0
	}
}

func UnicodeBytes2String(b []byte) string {
	decoder := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewDecoder()
	str, _ := decoder.Bytes(b)
	return string(str)
}

func Bytes2String(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
