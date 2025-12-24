package binom

import (
	"encoding/binary"
	"strconv"
)

type binomEvent [9]byte

func EventFromBytes(byteString string) binomEvent {
	var be binomEvent
	b := []byte(byteString)

	copy(be[:], b)

	return be
}

func Event(index int8, val int64) binomEvent {
	var be binomEvent
	be.setHeader("event", index)
	be.setValue(val)

	return be
}

func AddEvent(index int8, val int64) binomEvent {
	var be binomEvent
	be.setHeader("add_event", index)
	be.setValue(val)

	return be
}

func (ev binomEvent) Type() string {
	t, _ := ev.Header()
	return t
}

func (ev binomEvent) Index() int8 {
	_, i := ev.Header()
	return i
}

func (ev binomEvent) Name() string {
	t, i := ev.Header()
	return t + strconv.Itoa(int(i))
}

func (ev binomEvent) Value() int64 {
	return int64(binary.BigEndian.Uint64(ev[1:]))
}

func (ev binomEvent) URLParam() string {
	return ev.Name() + "=" + strconv.Itoa(int(ev.Value()))
}

func (ev binomEvent) Header() (string, int8) {
	t := "event"
	fb := int8(ev[0])

	if fb&0x1 == 0x1 {
		t = "add_event"
	}
	i := clearBit(fb, 7)

	return t, i
}

func (ev binomEvent) Bytes() []byte {
	var b []byte
	for _, v := range ev {
		b = append(b, v)
	}
	return []byte(b)
}

func setBit(n int8, pos uint) int8 {
	n |= (1 << pos)
	return n
}

func clearBit(n int8, pos uint) int8 {
	return n & ^(1 << pos)
}

func (ev *binomEvent) setHeader(t string, i int8) {
	fb := byte(i)
	if t == "add_event" {
		fb = byte(setBit(i, 7))
	}
	ev[0] = fb
}

func (ev *binomEvent) setValue(val int64) {
	binary.BigEndian.PutUint64(ev[1:], uint64(val))
}
