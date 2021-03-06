package Core

import (
	. "encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
	"strings"
)

var (
	BytesOrder = BigEndian
)

type Packet struct {
	Buffer []byte
	Index  int
}

func NewPacket() (p *Packet) {
	p = new(Packet)
	p.Index = 0
	p.Buffer = make([]byte, 1024)
	return p
}

func NewPacket2(size int) (p *Packet) {
	p = new(Packet)
	p.Index = 0
	p.Buffer = make([]byte, size)
	return p
}

func NewPacket3(buffer []byte) (p *Packet) {
	p = new(Packet)
	p.Buffer = make([]byte, len(buffer))
	p.Index = 0
	copy(p.Buffer, buffer)
	return p
}

func NewPacketRef(buffer []byte) (p *Packet) {
	p = new(Packet)
	p.Buffer = buffer
	p.Index = 0
	return p
}

func (p *Packet) String() string {
	pString := strings.Replace(string(p.Buffer), "\x07", " ", -1)
	return fmt.Sprintf("len(%d) : % #X\n %s", len(p.Buffer), p.Buffer, pString)
}

func (p *Packet) Clone() (pn *Packet) {
	pn = new(Packet)
	pn.Buffer = make([]byte, p.Index+2)
	copy(pn.Buffer, p.Buffer[:p.Index])
	pn.Index = p.Index
	return pn
}

func (p *Packet) WCheck(size int) {
	if p.Index+size > cap(p.Buffer) {
		p.Resize(p.Index + size + 10)
	}
}

func (p *Packet) RCheck(size int) bool {
	if p.Index+size > cap(p.Buffer) {
		return false
	}
	return true
}

func (p *Packet) Resize(newSize int) {

	if cap(p.Buffer) > newSize {
		p.Buffer = p.Buffer[:newSize]
		return
	}
	if newSize-cap(p.Buffer) > 1024 {
		p.Buffer = append(p.Buffer, make([]byte, newSize-cap(p.Buffer))...)
	} else {
		p.Buffer = append(p.Buffer, make([]byte, 1024)...)
	}

	p.Buffer = p.Buffer[:cap(p.Buffer)]

	//temp := make([]byte, newSize)
	//copy(temp, p.Buffer)
	//p.Buffer = temp
}

func (p *Packet) WSkip(times int) {
	p.WCheck(times)
	p.Index += times
}

func (p *Packet) RSkip(times int) {
	if !p.RCheck(times) {
		panic("Reading outside of the packet!")
	}
	p.Index += times
}

func (p *Packet) WriteByte(b byte) {
	p.WCheck(1)
	p.Buffer[p.Index] = b
	p.Index++
}

func (p *Packet) WriteLen() {
	t := p.Index
	p.Index = 1
	p.WriteUInt16(uint16(t - 3))
	p.Index = t
}

func (p *Packet) WriteInt16(value int16) {
	p.WCheck(2)
	BytesOrder.PutUint16(p.Buffer[p.Index:], uint16(value))
	p.Index += 2
}

func (p *Packet) WriteUInt16(value uint16) {
	p.WCheck(2)
	BytesOrder.PutUint16(p.Buffer[p.Index:], value)
	p.Index += 2
}

func (p *Packet) WriteInt32(value int32) {
	p.WCheck(4)
	BytesOrder.PutUint32(p.Buffer[p.Index:], uint32(value))
	p.Index += 4
}

func (p *Packet) WriteUInt32(value uint32) {
	p.WCheck(4)
	BytesOrder.PutUint32(p.Buffer[p.Index:], value)
	p.Index += 4
}

func (p *Packet) WriteInt64(value int64) {
	p.WCheck(8)
	BytesOrder.PutUint64(p.Buffer[p.Index:], uint64(value))
	p.Index += 8
}

func (p *Packet) WriteUInt64(value uint64) {
	p.WCheck(8)
	BytesOrder.PutUint64(p.Buffer[p.Index:], value)
	p.Index += 8
}

func (p *Packet) WriteFloat32(pValue float32) {
	p.WCheck(4)
	p.WriteUInt32(math.Float32bits(pValue))
}

func (p *Packet) WriteFloat64(pValue float64) {
	p.WCheck(8)
	p.WriteUInt64(math.Float64bits(pValue))
}

func (p *Packet) WriteString(pValue string, size int) {
	p.WCheck(size)
	copy(p.Buffer[p.Index:], pValue[:size])
	p.Index += size
}

func (p *Packet) WriteRawString(pValue string) {
	p.WriteString(pValue, len(pValue))
}

func (p *Packet) ReadBytes(size int) (pValue []byte) {
	if !p.RCheck(size) {
		panic("Reading outside of the packet!")
	}
	pValue = p.Buffer[p.Index : p.Index+size]
	p.Index += size
	return pValue
}

func (p *Packet) ReadByte() (pValue byte) {
	if !p.RCheck(1) {
		panic("Reading outside of the packet!")
	}
	pValue = byte(p.Buffer[p.Index])
	p.Index++
	return pValue
}

func (p *Packet) ReadInt16() (pValue int16) {
	if !p.RCheck(2) {
		panic("Reading outside of the packet!")
	}
	pValue = int16(BytesOrder.Uint16(p.Buffer[p.Index:]))
	p.Index += 2
	return pValue
}

func (p *Packet) ReadUInt16() (pValue uint16) {
	if !p.RCheck(2) {
		panic("Reading outside of the packet!")
	}
	pValue = BytesOrder.Uint16(p.Buffer[p.Index:])
	p.Index += 2
	return pValue
}

func (p *Packet) ReadInt32() (pValue int32) {
	if !p.RCheck(4) {
		panic("Reading outside of the packet!")
	}
	pValue = int32(BytesOrder.Uint32(p.Buffer[p.Index:]))
	p.Index += 4
	return pValue
}

func (p *Packet) ReadUInt32() (pValue uint32) {
	if !p.RCheck(4) {
		panic("Reading outside of the packet!")
	}
	pValue = BytesOrder.Uint32(p.Buffer[p.Index:])
	p.Index += 4
	return pValue
}

func (p *Packet) ReadInt64() (pValue int64) {
	if !p.RCheck(4) {
		panic("Reading outside of the packet!")
	}
	pValue = int64(BytesOrder.Uint64(p.Buffer[p.Index:]))
	p.Index += 8
	return pValue
}

func (p *Packet) ReadUInt64() (pValue uint64) {
	if !p.RCheck(8) {
		panic("Reading outside of the packet!")
	}
	pValue = BytesOrder.Uint64(p.Buffer[p.Index:])
	p.Index += 8
	return pValue
}

func (p *Packet) ReadFloat32() (pValue float32) {
	if !p.RCheck(4) {
		panic("Reading outside of the packet!")
	}
	pValue = math.Float32frombits(p.ReadUInt32())
	return pValue
}

func (p *Packet) ReadFloat64() (pValue float64) {
	if !p.RCheck(8) {
		panic("Reading outside of the packet!")
	}
	pValue = math.Float64frombits(p.ReadUInt64())
	return pValue
}

func (p *Packet) BasePacket() *Packet {
	return p
}

func (p *Packet) ReadString(size int) (pValue string) {
	if !p.RCheck(size) {
		panic("Reading outside of the packet!")
	}
	pValue = string(p.Buffer[p.Index : p.Index+size])
	p.Index += size
	return pValue
}

func (p *Packet) Read(b []byte) (n int, err error) {
	if len(b) == 0 {
		return 0, nil
	}
	if p.Buffer == nil {
		return 0, errors.New("nil buffer")
	}
	if p.Index >= cap(p.Buffer) {
		return 0, io.EOF
	}
	n = len(b)
	t := p.Index + n
	if t > cap(p.Buffer) {
		n = t - cap(p.Buffer)
	}
	copy(b, p.Buffer[p.Index:n])
	p.Index += n
	return n, nil
}

func (p *Packet) WriteBytes(bytes []byte) {
	p.WCheck(len(bytes))
	copy(p.Buffer[p.Index:], bytes)
	p.Index += len(bytes)
}

func (p *Packet) Write(bytes []byte) (n int, err error) {
	if len(bytes) == 0 {
		return 0, nil
	}
	if p.Buffer == nil {
		return 0, errors.New("nil buffer")
	}
	if p.Index >= cap(p.Buffer) {
		return 0, io.EOF
	}
	n = len(bytes)
	t := p.Index + n
	if t > cap(p.Buffer) {
		n = t - cap(p.Buffer)
	}
	copy(p.Buffer[p.Index:], bytes[:n])
	p.Index += n
	return n, nil
}

func (p *Packet) Bytes() []byte {
	return p.Buffer[:p.Index]
}
