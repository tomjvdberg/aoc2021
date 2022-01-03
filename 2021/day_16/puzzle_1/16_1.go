package main

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"time"
)

type transmission struct {
	data    []byte
	buffer  big.Int
	pointer int
}

func transmissionFromString(s string) *transmission {
	newT := &transmission{}
	newT.data, _ = hex.DecodeString(s)
	newT.renewBuffer()

	return newT
}
func (t *transmission) renewBuffer() {
	if len(t.data) == 0 {
		return
	}
	t.buffer = *big.NewInt(int64(t.data[0]))
	t.data = t.data[1:]
	t.pointer = 7

}
func (t *transmission) nextBit() (uint, bool) {
	b := t.buffer.Bit(t.pointer)
	t.pointer--
	if t.pointer < 0 {
		t.renewBuffer()
	}

	// if pointer is negative then the transmission if completed
	if t.pointer < 0 {
		return b, true
	}

	return b, false
}

const LITERAL = 4

const VERSION = 0
const TYPE_ID = 1
const LIT_PACK = 2
const LENGTH_TYPE_ID = 3
const TOTAL_LENGTH_OF_SUBPACKETS = 4
const COUNT_OF_SUBPACKETS = 5
const KNOWN_LENGTH_SUBPACKETS = 6
const KNOWN_COUNT_SUBPACKETS = 7

type segment interface {
	getSegmentType() int
	pushBit(b uint, versions *[]int64) (complete bool)
	getValue() []int64
	getBit(i int) uint
	getValueOfBits(start int, end int) int64
}
type knownLengthSegment struct {
	segmentType int
	bitCnt      int
	bits        big.Int
}

func (s *knownLengthSegment) getSegmentType() int {
	return s.segmentType
}
func (s *knownLengthSegment) pushBit(b uint, versions *[]int64) (complete bool) {
	s.bits.SetBit(&s.bits, s.bitCnt, b)
	s.bitCnt--

	return s.bitCnt == -1
}
func (s *knownLengthSegment) getValue() []int64 {
	return []int64{s.bits.Int64()}
}
func (s *knownLengthSegment) getBit(i int) uint {
	return s.bits.Bit(i)
}
func (s *knownLengthSegment) getValueOfBits(start int, end int) int64 {
	x := big.NewInt(0)
	for i := start; i < end; i++ {
		x.SetBit(x, i, s.bits.Bit(i))
	}

	return x.Int64()
}

type knownLengthSubpacketSegment struct {
	bitCnt    int64
	packets   []packet
	curPacket packet
}

func (s *knownLengthSubpacketSegment) getSegmentType() int {
	return KNOWN_LENGTH_SUBPACKETS
}
func (s *knownLengthSubpacketSegment) pushBit(b uint, versions *[]int64) (complete bool) {
	packetComplete := s.curPacket.pushBit(b, versions)
	if packetComplete {
		s.packets = append(s.packets, s.curPacket)
		s.curPacket = *newPacket()
	}
	s.bitCnt--

	return s.bitCnt == -1
}
func (s *knownLengthSubpacketSegment) getValue() []int64 {
	vs := []int64{}
	for _, p := range s.packets {
		vs = append(vs, p.valueOfLiteralsConcatenated())
	}

	return vs
}
func (s *knownLengthSubpacketSegment) getBit(i int) uint {
	return 0
}
func (s *knownLengthSubpacketSegment) getValueOfBits(start int, end int) int64 {
	return 0
}

func newKnownLengthSubpacketSegment(length int64) *knownLengthSubpacketSegment {
	pack := packet{}
	pack.curSegment = newVersionSegment()

	return &knownLengthSubpacketSegment{
		length - 1,
		[]packet{},
		pack,
	}
}

type countSubpacketSegment struct {
	packetCnt int64
	packets   []packet
	curPacket packet
}

func (s *countSubpacketSegment) getSegmentType() int {
	return KNOWN_COUNT_SUBPACKETS
}
func (s *countSubpacketSegment) pushBit(b uint, versions *[]int64) (complete bool) {
	packetComplete := s.curPacket.pushBit(b, versions)
	if packetComplete {
		s.packets = append(s.packets, s.curPacket)
		s.curPacket = *newPacket()
	}

	return int64(len(s.packets)) == s.packetCnt
}
func (s *countSubpacketSegment) getValue() []int64 {
	vs := []int64{}
	for _, p := range s.packets {
		vs = append(vs, p.valueOfLiteralsConcatenated())
	}

	return vs
}
func (s *countSubpacketSegment) getBit(i int) uint {
	return 0
}
func (s *countSubpacketSegment) getValueOfBits(start int, end int) int64 {
	return 0
}
func newCountSubpacketSegment(subPacketsExpected int64) *countSubpacketSegment {
	pack := packet{}
	pack.curSegment = newVersionSegment()

	return &countSubpacketSegment{
		subPacketsExpected,
		[]packet{},
		pack,
	}
}

type packet struct {
	version       int64
	typeId        int64
	literalValues []int64
	complete      bool
	curSegment    segment
	subpackets    []packet
}

func (p *packet) pushBit(b uint, versions *[]int64) bool {
	segmentIsComplete := p.curSegment.pushBit(b, versions)

	if !segmentIsComplete {
		return false
	}
	switch p.curSegment.getSegmentType() {
	case VERSION:
		p.version = p.curSegment.getValue()[0]
		*versions = append(*versions, p.version)
		p.curSegment = newTypeIdSegment()
	case TYPE_ID:
		p.typeId = p.curSegment.getValue()[0]
		if p.typeId == LITERAL {
			p.curSegment = newLiteralPacketSegment()
			break
		}
		// we're talking about an operator
		// so the length type id is next
		p.curSegment = newLengthTypeIdSegment()
	case LENGTH_TYPE_ID:
		if p.curSegment.getValue()[0] == 0 {
			// add the 15 bit length knownLengthSegment
			p.curSegment = newTotalLengthSubpacketsSegment()
			return false
		}
		p.curSegment = newCountOfSubpacketsSegment()
	case TOTAL_LENGTH_OF_SUBPACKETS:
		totalLengthOfSubpackets := p.curSegment.getValue()[0]
		p.curSegment = newKnownLengthSubpacketSegment(totalLengthOfSubpackets)
	case COUNT_OF_SUBPACKETS:
		countOfSubpackets := p.curSegment.getValue()[0]
		p.curSegment = newCountSubpacketSegment(countOfSubpackets)
	case KNOWN_LENGTH_SUBPACKETS:
		p.literalValues = append(p.literalValues, p.curSegment.getValue()...)
		return true
	case KNOWN_COUNT_SUBPACKETS:
		p.literalValues = append(p.literalValues, p.curSegment.getValue()...)
		return true
	case LIT_PACK:
		// check last 4 bits and add to literalValues
		p.literalValues = append(p.literalValues, p.curSegment.getValueOfBits(0, 4))
		// check first bit to see if another lit pack should be added
		if p.curSegment.getBit(4) == 1 {
			p.curSegment = newLiteralPacketSegment()
		} else {
			return true // the packet is complete
		}
	}

	return false
}

func (p *packet) valueOfLiteralsConcatenated() int64 {
	litPointer := (len(p.literalValues) * 4) - 1 // the value bit length
	x := big.NewInt(0)
	for _, p := range p.literalValues {
		pv := big.NewInt(p)
		for i := 3; i > -1; i-- {
			x.SetBit(x, litPointer, pv.Bit(i))
			litPointer--
		}
	}

	return x.Int64()
}

func newPacket() *packet {
	pack := packet{}
	pack.curSegment = newVersionSegment() // always start with this

	return &pack
}
func main() {
	start := time.Now()
	dataString := "6051639005B56008C1D9BB3CC9DAD5BE97A4A9104700AE76E672DC95AAE91425EF6AD8BA5591C00F92073004AC0171007E0BC248BE0008645982B1CA680A7A0CC60096802723C94C265E5B9699E7E94D6070C016958F99AC015100760B45884600087C6E88B091C014959C83E740440209FC89C2896A50765A59CE299F3640D300827902547661964D2239180393AF92A8B28F4401BCC8ED52C01591D7E9D2591D7E9D273005A5D127C99802C095B044D5A19A73DC0E9C553004F000DE953588129E372008F2C0169FDB44FA6C9219803E00085C378891F00010E8FF1AE398803D1BE25C743005A6477801F59CC4FA1F3989F420C0149ED9CF006A000084C5386D1F4401F87310E313804D33B4095AFBED32ABF2CA28007DC9D3D713300524BCA940097CA8A4AF9F4C00F9B6D00088654867A7BC8BCA4829402F9D6895B2E4DF7E373189D9BE6BF86B200B7E3C68021331CD4AE6639A974232008E663C3FE00A4E0949124ED69087A848002749002151561F45B3007218C7A8FE600FC228D50B8C01097EEDD7001CF9DE5C0E62DEB089805330ED30CD3C0D3A3F367A40147E8023221F221531C9681100C717002100B36002A19809D15003900892601F950073630024805F400150D400A70028C00F5002C00252600698400A700326C0E44590039687B313BF669F35C9EF974396EF0A647533F2011B340151007637C46860200D43085712A7E4FE60086003E5234B5A56129C91FC93F1802F12EC01292BD754BCED27B92BD754BCED27B100264C4C40109D578CA600AC9AB5802B238E67495391D5CFC402E8B325C1E86F266F250B77ECC600BE006EE00085C7E8DF044001088E31420BCB08A003A72BF87D7A36C994CE76545030047801539F649BF4DEA52CBCA00B4EF3DE9B9CFEE379F14608"
	trans := transmissionFromString(dataString)

	versions := []int64{}
	packets := []packet{}
	pack := newPacket()

	for {
		nb, transmissionFinished := trans.nextBit()
		packComplete := (*pack).pushBit(nb, &versions)

		if packComplete {
			packets = append(packets, *pack)
			pack = newPacket()
		}

		if transmissionFinished {
			break
		}
	}
	versionSum := int64(0)
	for _, version := range versions {
		versionSum += version
	}
	fmt.Println(versionSum)
	fmt.Println("End", time.Since(start))
}

func createKnownLengthSegment(len int, segmentType int) *knownLengthSegment {
	return &knownLengthSegment{
		segmentType,
		len - 1,
		*big.NewInt(0),
	}
}

func newVersionSegment() *knownLengthSegment {
	return createKnownLengthSegment(3, VERSION)
}
func newTypeIdSegment() *knownLengthSegment {
	return createKnownLengthSegment(3, TYPE_ID)
}
func newLiteralPacketSegment() *knownLengthSegment {
	return createKnownLengthSegment(5, LIT_PACK)
}
func newLengthTypeIdSegment() *knownLengthSegment {
	return createKnownLengthSegment(1, LENGTH_TYPE_ID)
}
func newTotalLengthSubpacketsSegment() *knownLengthSegment {
	return createKnownLengthSegment(15, TOTAL_LENGTH_OF_SUBPACKETS)
}
func newCountOfSubpacketsSegment() *knownLengthSegment {
	return createKnownLengthSegment(11, COUNT_OF_SUBPACKETS)
}
