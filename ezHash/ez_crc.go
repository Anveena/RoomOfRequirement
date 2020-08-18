package ezHash

type Parameters struct {
	Width      uint
	Polynomial uint64
	ReflectIn  bool
	ReflectOut bool
	Init       uint64
	FinalXor   uint64
}

func reflect(in uint64, count uint) uint64 {
	rs := in
	for idx := uint(0); idx < count; idx++ {
		srcbit := uint64(1) << idx
		dstbit := uint64(1) << (count - idx - 1)
		if (in & srcbit) != 0 {
			rs |= dstbit
		} else {
			rs = rs & (^dstbit)
		}
	}
	return rs
}
func CalculateCRC(crcParams *Parameters, data []byte) uint64 {
	curValue := crcParams.Init
	topBit := uint64(1) << (crcParams.Width - 1)
	mask := (topBit << 1) - 1
	for i := 0; i < len(data); i++ {
		var curByte = uint64(data[i]) & 0x00FF
		if crcParams.ReflectIn {
			curByte = reflect(curByte, 8)
		}
		for j := uint64(0x0080); j != 0; j >>= 1 {
			bit := curValue & topBit
			curValue <<= 1
			if (curByte & j) != 0 {
				bit = bit ^ topBit
			}
			if bit != 0 {
				curValue = curValue ^ crcParams.Polynomial
			}
		}
	}
	if crcParams.ReflectOut {
		curValue = reflect(curValue, crcParams.Width)
	}
	curValue = curValue ^ crcParams.FinalXor
	return curValue & mask
}

type Table struct {
	crcParams Parameters
	crcTable  []uint64
	mask      uint64
	initValue uint64
}

func NewTable(crcParams *Parameters) *Table {
	ret := &Table{crcParams: *crcParams}
	ret.mask = (uint64(1) << crcParams.Width) - 1
	ret.crcTable = make([]uint64, 256, 256)
	ret.initValue = crcParams.Init
	if crcParams.ReflectIn {
		ret.initValue = reflect(crcParams.Init, crcParams.Width)
	}
	tmp := make([]byte, 1, 1)
	tableParams := *crcParams
	tableParams.Init = 0
	tableParams.ReflectOut = tableParams.ReflectIn
	tableParams.FinalXor = 0
	for i := 0; i < 256; i++ {
		tmp[0] = byte(i)
		ret.crcTable[i] = CalculateCRC(&tableParams, tmp)
	}
	return ret
}
func (t *Table) InitCrc() uint64 {
	return t.initValue
}
func (t *Table) UpdateCrc(curValue uint64, p []byte) uint64 {
	if t.crcParams.ReflectIn {
		for _, v := range p {
			curValue = t.crcTable[(byte(curValue)^v)&0xff] ^ (curValue >> 8)
		}
	} else if t.crcParams.Width < 8 {
		for _, v := range p {
			curValue = t.crcTable[((((byte)(curValue<<(8-t.crcParams.Width)))^v)&0xff)] ^ (curValue << 8)
		}
	} else {
		for _, v := range p {
			curValue = t.crcTable[((byte(curValue>>(t.crcParams.Width-8))^v)&0xff)] ^ (curValue << 8)
		}
	}
	return curValue
}
func (t *Table) CRC(curValue uint64) uint64 {
	ret := curValue

	if t.crcParams.ReflectOut != t.crcParams.ReflectIn {
		ret = reflect(ret, t.crcParams.Width)
	}
	return (ret ^ t.crcParams.FinalXor) & t.mask
}
func (t *Table) CalculateCRC(data []byte) uint64 {
	crc := t.InitCrc()
	crc = t.UpdateCrc(crc, data)
	return t.CRC(crc)
}

type Hash struct {
	table    *Table
	curValue uint64
	size     uint
}

func (h *Hash) Size() int      { return int(h.size) }
func (h *Hash) BlockSize() int { return 1 }
func (h *Hash) Reset() {
	h.curValue = h.table.InitCrc()
}
func (h *Hash) Sum(in []byte) []byte {
	s := h.CRC()
	for i := h.size; i > 0; {
		i--
		in = append(in, byte(s>>(8*i)))
	}
	return in
}
func (h *Hash) Write(p []byte) (n int, err error) {
	h.Update(p)
	return len(p), nil
}
func (h *Hash) Update(p []byte) {
	h.curValue = h.table.UpdateCrc(h.curValue, p)
}
func (h *Hash) CRC() uint64 {
	return h.table.CRC(h.curValue)
}
func (h *Hash) CalculateCRC(data []byte) uint64 {
	return h.table.CalculateCRC(data)
}
func NewHashWithTable(table *Table) *Hash {
	ret := &Hash{table: table}
	ret.size = (table.crcParams.Width + 7) / 8
	ret.Reset()
	return ret
}
func NewHash(crcParams *Parameters) *Hash {
	return NewHashWithTable(NewTable(crcParams))
}
func (h *Hash) Table() *Table {
	return h.table
}
