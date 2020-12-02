package rcnb

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEncode(t *testing.T) {
	d, err := Encode([]byte("rcnb"))
	assert.Equal(t, nil, err)
	assert.Equal(t, "ɌcńƁȓČņÞ", d)
	d, err = Encode([]byte("The world is based on RC thus RCNB"))
	assert.Equal(t, nil, err)
	assert.Equal(t, "ȐčnÞȒċƝÞɌČǹƄɌcŅƀȒĉȠƀȒȼȵƄŕƇŅbȒCŅßȒČnƅŕƇņƁȓċȵƀȐĉņþŕƇņÞȒȻŅBɌCňƀȐĉņþƦȻƝƃ", d)
}

func TestDecode(t *testing.T) {
	d, err := Decode("ɌcńƁȓČņÞ")
	assert.Equal(t, nil, err)
	assert.Equal(t, []byte("rcnb"), d)
	d, err = Decode("ȐčnÞȒċƝÞɌČǹƄɌcŅƀȒĉȠƀȒȼȵƄŕƇŅbȒCŅßȒČnƅŕƇņƁȓċȵƀȐĉņþŕƇņÞȒȻŅBɌCňƀȐĉņþƦȻƝƃ")
	assert.Equal(t, nil, err)
	assert.Equal(t, []byte("The world is based on RC thus RCNB"), d)
}
