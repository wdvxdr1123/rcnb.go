package rcnb

import (
	"github.com/pkg/errors"
	"math"
	"strings"
)

var (
	rString = "rRŔŕŖŗŘřƦȐȑȒȓɌɍ"
	cString = "cCĆćĈĉĊċČčƇƈÇȻȼ"
	nString = "nNŃńŅņŇňƝƞÑǸǹȠȵ"
	bString = "bBƀƁƃƄƅßÞþ"

	rMap = map[rune]int{}
	cMap = map[rune]int{}
	nMap = map[rune]int{}
	bMap = map[rune]int{}

	rChar = []rune(rString)
	cChar = []rune(cString)
	nChar = []rune(nString)
	bChar = []rune(bString)

	NotEnoughNB  = errors.New("not enough NB")
	RCNBOverflow = errors.New("rcnb overflow")
	LengthNotNB  = errors.New("length not NB")
)

func init() {
	for i, v := range rChar {
		rMap[v] = i
	}
	for i, v := range cChar {
		cMap[v] = i
	}
	for i, v := range nChar {
		nMap[v] = i
	}
	for i, v := range bChar {
		bMap[v] = i
	}
}

const (
	rLen   = 15
	cLen   = 15
	nLen   = 15
	bLen   = 10
	rcLen  = rLen * cLen
	nbLen  = nLen * bLen
	cnbLen = cLen * nbLen
)

func div(a int, b int) uint16 {
	return uint16(math.Floor(float64(a) / float64(b)))
}

func EncodeByte(value uint16) (string, error) {
	if value > 0xFF {
		return "", RCNBOverflow
	}
	var builder = strings.Builder{}
	if value > 0x7f {
		value = value & 0x7F
		builder.WriteRune(nChar[div(int(value), bLen)])
		builder.WriteRune(bChar[value%bLen])
	} else {
		builder.WriteRune(rChar[value/cLen])
		builder.WriteRune(cChar[value%cLen])
	}
	return builder.String(), nil
}

func EncodeShort(value int) (string, error) {
	if value > 0xFFFF {
		return "", RCNBOverflow
	}
	var builder = strings.Builder{}
	if value > 0x7FFF {
		value = value & 0x7FFF
		builder.WriteRune(nChar[div(value%nbLen, bLen)])
		builder.WriteRune(bChar[value%bLen])
		builder.WriteRune(rChar[div(value, cnbLen)])
		builder.WriteRune(cChar[div(value%cnbLen, nbLen)])
	} else {
		builder.WriteRune(rChar[div(value, cnbLen)])
		builder.WriteRune(cChar[div(value%cnbLen, nbLen)])
		builder.WriteRune(nChar[div(value%nbLen, bLen)])
		builder.WriteRune(bChar[value%bLen])
	}
	return builder.String(), nil
}

func Encode(bytes []byte) (string, error) {
	var (
		builder = strings.Builder{}
		length  = len(bytes)
	)
	for i := 0; i < (length >> 1); i++ {
		s, err := EncodeShort((int(bytes[i*2]) << 8) | int(bytes[i*2+1]))
		if err != nil {
			return "", err
		}
		builder.WriteString(s)
	}
	if length&1 == 1 {
		s, err := EncodeByte(uint16(bytes[length-1]))
		if err != nil {
			return "", err
		}
		builder.WriteString(s)
	}
	return builder.String(), nil
}

func load(m map[rune]int, r rune) int {
	if v, ok := m[r]; ok {
		return v
	}
	return -1
}

func DecodeByte(c []rune) (byte, error) {
	var (
		nb   = false
		idx1 = load(rMap, c[0])
		idx2 = load(cMap, c[1])
	)
	if len(c) != 2 {
		return 0, LengthNotNB
	}
	if idx1 < 0 || idx2 < 0 {
		idx1 = load(nMap, c[0])
		idx2 = load(bMap, c[1])
		nb = true
	}
	if idx1 < 0 || idx2 < 0 {
		return 0, NotEnoughNB
	}
	if nb {
		return byte((idx1*bLen + idx2) | 0x80), nil
	} else {
		return byte(idx1*cLen + idx2), nil
	}
}

func DecodeShort(c []rune) (int, error) {
	var idx [4]int
	if len(c) != 4 {
		return 0, LengthNotNB
	}
	reverse := load(rMap, c[0]) < 0
	if reverse {
		idx = [4]int{
			load(rMap, c[2]),
			load(cMap, c[3]),
			load(nMap, c[0]),
			load(bMap, c[1]),
		}
	} else {
		idx = [4]int{
			load(rMap, c[0]),
			load(cMap, c[1]),
			load(nMap, c[2]),
			load(bMap, c[3]),
		}
	}
	if idx[0] < 0 || idx[1] < 0 || idx[2] < 0 || idx[3] < 0 {
		return 0, NotEnoughNB
	}
	result := idx[0]*cnbLen + idx[1]*nbLen + idx[2]*bLen + idx[3]
	if result > 0x7FFF {
		return 0, RCNBOverflow
	}
	if reverse {
		result |= 0x8000
	}
	return result, nil
}

func Decode(s string) ([]byte, error) {
	var (
		c      = []rune(s)
		length = len(c)

		ret []byte
	)
	if length&1 == 1 {
		return nil, LengthNotNB
	}
	for i := 0; i < (length >> 2); i++ {
		s, err := DecodeShort(c[i*4 : i*4+4])
		if err != nil {
			return nil, err
		}
		ret = append(ret, byte(s>>8), byte(s&0xFF))
	}
	if length&2 == 2 {
		s, err := DecodeByte(c[length-2:])
		if err != nil {
			return nil, err
		}
		ret = append(ret, s)
	}
	return ret, nil
}
