package rdb

import ("encoding/binary")

type Writer struct {
	data []byte
}

func NewWriter() *Writer {
	return &Writer{
		data: make([]byte, 0),
	}
}

func (w *Writer)WriteByte(b byte) {
	w.data = append(w.data, b)
}

func (w *Writer)WriteBytes(data []byte) {
	w.data = append(w.data, data...)
}

func (w *Writer) WriteLength(length uint32) {
	switch {
	case length <= 63:
		w.WriteByte(byte(length))

	case length <= 16383:
		first := byte(0x40 | byte(length>>8))
		second := byte(length)

		w.WriteByte(first)
		w.WriteByte(second)

	default:
		w.WriteByte(0x80)

		buf := make([]byte, 4)
		binary.BigEndian.PutUint32(buf, length)

		w.WriteBytes(buf)
		}
}

func (w *Writer) WriteString(value string) {
	data := []byte(value)

	w.WriteLength(uint32(len(data)))
	w.WriteBytes(data)
}

func (w *Writer) Bytes() []byte {
	return w.data
}