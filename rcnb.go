package rcnb

import (
	"github.com/pkg/errors"
	"math"
	"strings"
)

var (
	rChar = []rune("rRŔŕŖŗŘřƦȐȑȒȓɌɍ")
	cChar = []rune("cCĆćĈĉĊċČčƇƈÇȻȼ")
	nChar = []rune("nNŃńŅņŇňƝƞÑǸǹȠȵ")
	bChar = []rune("bBƀƁƃƄƅßÞþ")
)

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
		return "", errors.New("rcnb overflow")
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
		return "", errors.New("rcnb overflow")
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
	var builder = strings.Builder{}
	var length = len(bytes)
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
